package services

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sparklassignment/backend/internal/dto"
	"sparklassignment/backend/internal/models"
)

type StudentService struct {
	db *gorm.DB
}

func NewStudentService(db *gorm.DB) *StudentService {
	return &StudentService{
		db: db,
	}
}

func (s *StudentService) GetQuizzes(studentID uint) (*dto.StudentQuizResponse, error) {
	var quizzes []models.Quiz
	if err := s.db.Order("created_at desc").Find(&quizzes).Error; err != nil {
		return nil, err
	}

	quizMap := make(map[uint]models.Quiz, len(quizzes))
	for _, quiz := range quizzes {
		quizMap[quiz.QuizID] = quiz
	}

	var attempts []models.QuizAttempt
	if err := s.db.Where("student_id = ?", studentID).Order("updated_at desc").Find(&attempts).Error; err != nil {
		return nil, err
	}

	for index := range attempts {
		quiz, exists := quizMap[attempts[index].QuizID]
		if !exists {
			continue
		}

		if attemptExpired(attempts[index], quiz) {
			if err := finalizeAttempt(s.db, &attempts[index]); err != nil {
				return nil, err
			}
		}
	}

	latestAttemptByQuiz := make(map[uint]models.QuizAttempt)
	for _, attempt := range attempts {
		if _, exists := latestAttemptByQuiz[attempt.QuizID]; !exists {
			latestAttemptByQuiz[attempt.QuizID] = attempt
		}
	}

	response := make([]dto.StudentQuizItemResponse, 0, len(quizzes))
	for _, quiz := range quizzes {
		item := dto.StudentQuizItemResponse{
			ID:              quiz.QuizID,
			Category:        quiz.Category,
			Title:           quiz.Title,
			QuestionCount:   quiz.QuestionCount,
			TotalMarks:      quiz.TotalMarks,
			DurationMinutes: quiz.DurationMinutes,
			Status:          "not_started",
			Action:          "start",
		}

		if attempt, exists := latestAttemptByQuiz[quiz.QuizID]; exists {
			attemptID := attempt.QuizAttemptID
			switch attempt.Status {
			case "in_progress":
				item.Status = "in_progress"
				item.Action = "resume"
				item.AttemptID = &attemptID
			case "completed":
				item.Status = "completed"
				item.Action = "view_score"
				item.AttemptID = &attemptID
			}
		}

		response = append(response, item)
	}

	return &dto.StudentQuizResponse{Quizzes: response}, nil
}

func (s *StudentService) StartQuiz(studentID, quizID uint) (*dto.StartQuizResponse, string, error) {
	var quiz models.Quiz
	if err := s.db.First(&quiz, quizID).Error; err != nil {
		return nil, "", err
	}

	var inProgressAttempt models.QuizAttempt
	err := s.db.Where("quiz_id = ? AND student_id = ? AND status = ?", quizID, studentID, "in_progress").First(&inProgressAttempt).Error
	if err == nil {
		if attemptExpired(inProgressAttempt, quiz) {
			if err := finalizeAttempt(s.db, &inProgressAttempt); err != nil {
				return nil, "", err
			}

			return &dto.StartQuizResponse{
				Message:   "quiz already completed",
				AttemptID: inProgressAttempt.QuizAttemptID,
				Status:    "completed",
			}, "completed", nil
		}

		return &dto.StartQuizResponse{
			Message:   "quiz already in progress",
			AttemptID: inProgressAttempt.QuizAttemptID,
			Status:    "in_progress",
		}, "resume", nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	var completedAttempt models.QuizAttempt
	err = s.db.Where("quiz_id = ? AND student_id = ? AND status = ?", quizID, studentID, "completed").Order("updated_at desc").First(&completedAttempt).Error
	if err == nil {
		return &dto.StartQuizResponse{
			Message:   "quiz already completed",
			AttemptID: completedAttempt.QuizAttemptID,
			Status:    "completed",
		}, "completed", nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	attempt := models.QuizAttempt{
		QuizID:    quizID,
		StudentID: studentID,
		Status:    "in_progress",
		StartedAt: time.Now(),
		Score:     0,
	}

	if err := s.db.Create(&attempt).Error; err != nil {
		existingAttempt, fetchErr := s.findLatestAttemptForQuiz(studentID, quizID)
		if fetchErr == nil {
			if existingAttempt.Status == "in_progress" {
				if attemptExpired(*existingAttempt, quiz) {
					if err := finalizeAttempt(s.db, existingAttempt); err != nil {
						return nil, "", err
					}

					return &dto.StartQuizResponse{
						Message:   "quiz already completed",
						AttemptID: existingAttempt.QuizAttemptID,
						Status:    "completed",
					}, "completed", nil
				}

				return &dto.StartQuizResponse{
					Message:   "quiz already in progress",
					AttemptID: existingAttempt.QuizAttemptID,
					Status:    "in_progress",
				}, "resume", nil
			}

			if existingAttempt.Status == "completed" {
				return &dto.StartQuizResponse{
					Message:   "quiz already completed",
					AttemptID: existingAttempt.QuizAttemptID,
					Status:    "completed",
				}, "completed", nil
			}
		}

		return nil, "", err
	}

	return &dto.StartQuizResponse{
		Message:   "quiz started successfully",
		AttemptID: attempt.QuizAttemptID,
		Status:    attempt.Status,
	}, "started", nil
}

func (s *StudentService) GetAttempt(studentID, attemptID uint) (*dto.GetAttemptResponse, error) {
	attempt, err := loadAttemptForStudent(s.db, attemptID, studentID)
	if err != nil {
		return nil, err
	}

	bundle, err := loadAttemptBundle(s.db, *attempt)
	if err != nil {
		return nil, err
	}

	if attemptExpired(bundle.Attempt, bundle.Quiz) {
		if err := finalizeAttempt(s.db, attempt); err != nil {
			return nil, err
		}

		bundle, err = loadAttemptBundle(s.db, *attempt)
		if err != nil {
			return nil, err
		}
	}

	questions, err := buildAttemptQuestions(bundle)
	if err != nil {
		return nil, err
	}

	return &dto.GetAttemptResponse{
		AttemptID:        bundle.Attempt.QuizAttemptID,
		Status:           bundle.Attempt.Status,
		RemainingSeconds: remainingSeconds(bundle.Attempt, bundle.Quiz),
		Quiz:             toAttemptQuizResponse(bundle.Quiz),
		Questions:        questions,
	}, nil
}

func (s *StudentService) SaveAnswer(studentID, attemptID uint, req dto.SaveAnswerRequest) (*dto.SaveAnswerResponse, error) {
	attempt, err := loadAttemptForStudent(s.db, attemptID, studentID)
	if err != nil {
		return nil, err
	}

	bundle, err := loadAttemptBundle(s.db, *attempt)
	if err != nil {
		return nil, err
	}

	if attemptExpired(bundle.Attempt, bundle.Quiz) {
		if err := finalizeAttempt(s.db, attempt); err != nil {
			return nil, err
		}

		return nil, ConflictError{Message: "attempt already finished"}
	}

	if bundle.Attempt.Status == "completed" {
		return nil, ConflictError{Message: "attempt already finished"}
	}

	quizQuestion, err := findQuizQuestion(bundle.QuizQuestions, req.QuizQuestionID)
	if err != nil {
		return nil, ValidationError{Message: "quiz question does not belong to this attempt"}
	}

	question, exists := bundle.QuestionsByID[quizQuestion.QuestionID]
	if !exists {
		return nil, errors.New("question not found")
	}

	allowedOptions, err := ParseJSONStringArray(question.Options)
	if err != nil {
		return nil, err
	}

	if !validateChosenOptions(req.ChosenOptions, allowedOptions) {
		return nil, ValidationError{Message: "chosen_options contains invalid values"}
	}

	correctOptions, err := ParseJSONStringArray(question.CorrectOptions)
	if err != nil {
		return nil, err
	}

	rawChosenOptions, err := json.Marshal(req.ChosenOptions)
	if err != nil {
		return nil, err
	}

	awardedMarks := CalculateAwardedMarks(req.ChosenOptions, correctOptions, quizQuestion.Marks)
	answer := models.AttemptAnswer{
		AttemptID:      attemptID,
		QuizQuestionID: req.QuizQuestionID,
		ChosenOptions:  rawChosenOptions,
		AwardedMarks:   awardedMarks,
	}

	if err := s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "attempt_id"},
			{Name: "quiz_question_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"chosen_options": rawChosenOptions,
			"awarded_marks":  awardedMarks,
			"updated_at":     time.Now(),
		}),
	}).Create(&answer).Error; err != nil {
		return nil, err
	}

	return &dto.SaveAnswerResponse{
		Message:        "answer saved successfully",
		AttemptID:      attemptID,
		QuizQuestionID: req.QuizQuestionID,
	}, nil
}

func (s *StudentService) findLatestAttemptForQuiz(studentID, quizID uint) (*models.QuizAttempt, error) {
	var attempt models.QuizAttempt
	if err := s.db.
		Where("quiz_id = ? AND student_id = ?", quizID, studentID).
		Order("updated_at desc").
		First(&attempt).Error; err != nil {
		return nil, err
	}

	return &attempt, nil
}

func (s *StudentService) FinishAttempt(studentID, attemptID uint) (*dto.FinishAttemptResponse, error) {
	attempt, err := loadAttemptForStudent(s.db, attemptID, studentID)
	if err != nil {
		return nil, err
	}

	if attempt.Status == "in_progress" {
		if err := finalizeAttempt(s.db, attempt); err != nil {
			return nil, err
		}
	}

	return &dto.FinishAttemptResponse{
		Message:     "attempt finished successfully",
		AttemptID:   attempt.QuizAttemptID,
		Status:      attempt.Status,
		Score:       attempt.Score,
		SubmittedAt: attempt.SubmittedAt,
	}, nil
}

func (s *StudentService) GetReport(studentID, attemptID uint) (*dto.AttemptReportResponse, error) {
	attempt, err := loadAttemptForStudent(s.db, attemptID, studentID)
	if err != nil {
		return nil, err
	}

	bundle, err := loadAttemptBundle(s.db, *attempt)
	if err != nil {
		return nil, err
	}

	if attemptExpired(bundle.Attempt, bundle.Quiz) {
		if err := finalizeAttempt(s.db, attempt); err != nil {
			return nil, err
		}

		bundle, err = loadAttemptBundle(s.db, *attempt)
		if err != nil {
			return nil, err
		}
	}

	if bundle.Attempt.Status != "completed" {
		return nil, ConflictError{Message: "attempt is still in progress"}
	}

	return buildAttemptReport(bundle)
}
