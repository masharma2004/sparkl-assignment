package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sparklassignment/backend/internal/dto"
	"sparklassignment/backend/internal/models"
)

type attemptBundle struct {
	Attempt               models.QuizAttempt
	Quiz                  models.Quiz
	QuizQuestions         []models.QuizQuestion
	QuestionsByID         map[uint]models.Question
	AnswersByQuizQuestion map[uint]models.AttemptAnswer
}

func loadAttemptForStudent(database *gorm.DB, attemptID, studentID uint) (*models.QuizAttempt, error) {
	var attempt models.QuizAttempt
	if err := database.Where("quiz_attempt_id = ? AND student_id = ?", attemptID, studentID).First(&attempt).Error; err != nil {
		return nil, err
	}

	return &attempt, nil
}

func loadAttemptBundle(database *gorm.DB, attempt models.QuizAttempt) (*attemptBundle, error) {
	var quiz models.Quiz
	if err := database.First(&quiz, attempt.QuizID).Error; err != nil {
		return nil, fmt.Errorf("fetch quiz: %w", err)
	}

	var quizQuestions []models.QuizQuestion
	if err := database.Where("quiz_id = ?", quiz.QuizID).Order("sequence_number asc").Find(&quizQuestions).Error; err != nil {
		return nil, fmt.Errorf("fetch quiz questions: %w", err)
	}

	questionIDs := make([]uint, 0, len(quizQuestions))
	for _, quizQuestion := range quizQuestions {
		questionIDs = append(questionIDs, quizQuestion.QuestionID)
	}

	questionsByID := make(map[uint]models.Question, len(questionIDs))
	if len(questionIDs) > 0 {
		var questions []models.Question
		if err := database.Where("question_id IN ?", questionIDs).Find(&questions).Error; err != nil {
			return nil, fmt.Errorf("fetch questions: %w", err)
		}

		for _, question := range questions {
			questionsByID[question.QuestionID] = question
		}
	}

	var answers []models.AttemptAnswer
	if err := database.Where("attempt_id = ?", attempt.QuizAttemptID).Find(&answers).Error; err != nil {
		return nil, fmt.Errorf("fetch answers: %w", err)
	}

	answersByQuizQuestion := make(map[uint]models.AttemptAnswer, len(answers))
	for _, answer := range answers {
		answersByQuizQuestion[answer.QuizQuestionID] = answer
	}

	return &attemptBundle{
		Attempt:               attempt,
		Quiz:                  quiz,
		QuizQuestions:         quizQuestions,
		QuestionsByID:         questionsByID,
		AnswersByQuizQuestion: answersByQuizQuestion,
	}, nil
}

func toAttemptQuizResponse(quiz models.Quiz) dto.AttemptQuizResponse {
	return dto.AttemptQuizResponse{
		ID:              quiz.QuizID,
		Category:        quiz.Category,
		Title:           quiz.Title,
		DurationMinutes: quiz.DurationMinutes,
		QuestionCount:   quiz.QuestionCount,
		TotalMarks:      quiz.TotalMarks,
	}
}

func remainingSeconds(attempt models.QuizAttempt, quiz models.Quiz) int64 {
	if attempt.Status != "in_progress" {
		return 0
	}

	endTime := attempt.StartedAt.Add(time.Duration(quiz.DurationMinutes) * time.Minute)
	remaining := int64(time.Until(endTime).Seconds())
	if remaining < 0 {
		return 0
	}

	return remaining
}

func attemptExpired(attempt models.QuizAttempt, quiz models.Quiz) bool {
	if attempt.Status != "in_progress" {
		return false
	}

	endTime := attempt.StartedAt.Add(time.Duration(quiz.DurationMinutes) * time.Minute)
	return !time.Now().Before(endTime)
}

func buildAttemptQuestions(bundle *attemptBundle) ([]dto.AttemptQuestionResponse, error) {
	items := make([]dto.AttemptQuestionResponse, 0, len(bundle.QuizQuestions))

	for _, quizQuestion := range bundle.QuizQuestions {
		question, exists := bundle.QuestionsByID[quizQuestion.QuestionID]
		if !exists {
			return nil, fmt.Errorf("question %d not found in bundle", quizQuestion.QuestionID)
		}

		options, err := ParseJSONStringArray(question.Options)
		if err != nil {
			return nil, fmt.Errorf("parse question options: %w", err)
		}

		chosenOptions := []string{}
		if answer, exists := bundle.AnswersByQuizQuestion[quizQuestion.QuizQuestionID]; exists {
			chosenOptions, err = ParseJSONStringArray(answer.ChosenOptions)
			if err != nil {
				return nil, fmt.Errorf("parse chosen options: %w", err)
			}
		}

		items = append(items, dto.AttemptQuestionResponse{
			QuizQuestionID: quizQuestion.QuizQuestionID,
			QuestionID:     question.QuestionID,
			SequenceNumber: quizQuestion.SequenceNumber,
			Marks:          quizQuestion.Marks,
			Prompt:         question.Prompt,
			Options:        options,
			ChosenOptions:  chosenOptions,
		})
	}

	return items, nil
}

func buildAttemptReport(bundle *attemptBundle) (*dto.AttemptReportResponse, error) {
	reportQuestions := make([]dto.AttemptReportQuestionResponse, 0, len(bundle.QuizQuestions))
	totalScore := 0

	for _, quizQuestion := range bundle.QuizQuestions {
		question, exists := bundle.QuestionsByID[quizQuestion.QuestionID]
		if !exists {
			return nil, fmt.Errorf("question %d not found in bundle", quizQuestion.QuestionID)
		}

		options, err := ParseJSONStringArray(question.Options)
		if err != nil {
			return nil, fmt.Errorf("parse options: %w", err)
		}

		correctOptions, err := ParseJSONStringArray(question.CorrectOptions)
		if err != nil {
			return nil, fmt.Errorf("parse correct options: %w", err)
		}

		chosenOptions := []string{}
		if answer, exists := bundle.AnswersByQuizQuestion[quizQuestion.QuizQuestionID]; exists {
			chosenOptions, err = ParseJSONStringArray(answer.ChosenOptions)
			if err != nil {
				return nil, fmt.Errorf("parse chosen options: %w", err)
			}
		}

		awardedMarks := CalculateAwardedMarks(chosenOptions, correctOptions, quizQuestion.Marks)
		totalScore += awardedMarks

		reportQuestions = append(reportQuestions, dto.AttemptReportQuestionResponse{
			QuizQuestionID: quizQuestion.QuizQuestionID,
			QuestionID:     question.QuestionID,
			SequenceNumber: quizQuestion.SequenceNumber,
			Prompt:         question.Prompt,
			Options:        options,
			ChosenOptions:  chosenOptions,
			CorrectOptions: correctOptions,
			Solution:       question.Solution,
			AwardedMarks:   awardedMarks,
			MaxMarks:       quizQuestion.Marks,
			IsCorrect:      awardedMarks == quizQuestion.Marks,
		})
	}

	return &dto.AttemptReportResponse{
		AttemptID:   bundle.Attempt.QuizAttemptID,
		Status:      bundle.Attempt.Status,
		Score:       totalScore,
		SubmittedAt: bundle.Attempt.SubmittedAt,
		Quiz:        toAttemptQuizResponse(bundle.Quiz),
		Questions:   reportQuestions,
	}, nil
}

func finalizeAttempt(database *gorm.DB, attempt *models.QuizAttempt) error {
	if attempt.Status == "completed" {
		return nil
	}

	bundle, err := loadAttemptBundle(database, *attempt)
	if err != nil {
		return err
	}

	report, err := buildAttemptReport(bundle)
	if err != nil {
		return err
	}

	tx := database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, question := range report.Questions {
		if answer, exists := bundle.AnswersByQuizQuestion[question.QuizQuestionID]; exists {
			if err := tx.Model(&models.AttemptAnswer{}).
				Where("attempt_answer_id = ?", answer.AttemptAnswerID).
				Update("awarded_marks", question.AwardedMarks).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("update awarded marks: %w", err)
			}
		}
	}

	submittedAt := time.Now()
	if err := tx.Model(&models.QuizAttempt{}).
		Where("quiz_attempt_id = ?", attempt.QuizAttemptID).
		Updates(map[string]any{
			"status":       "completed",
			"score":        report.Score,
			"submitted_at": submittedAt,
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("update attempt: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit attempt update: %w", err)
	}

	attempt.Status = "completed"
	attempt.Score = report.Score
	attempt.SubmittedAt = &submittedAt

	return nil
}

func validateChosenOptions(chosenOptions, allowedOptions []string) bool {
	allowedSet := make(map[string]struct{}, len(allowedOptions))
	for _, option := range allowedOptions {
		allowedSet[option] = struct{}{}
	}

	seenChosenOptions := make(map[string]struct{}, len(chosenOptions))
	for _, chosenOption := range chosenOptions {
		if _, exists := allowedSet[chosenOption]; !exists {
			return false
		}

		if _, exists := seenChosenOptions[chosenOption]; exists {
			return false
		}

		seenChosenOptions[chosenOption] = struct{}{}
	}

	return true
}

func findQuizQuestion(quizQuestions []models.QuizQuestion, quizQuestionID uint) (*models.QuizQuestion, error) {
	for _, quizQuestion := range quizQuestions {
		if quizQuestion.QuizQuestionID == quizQuestionID {
			return &quizQuestion, nil
		}
	}

	return nil, errors.New("quiz question not found")
}
