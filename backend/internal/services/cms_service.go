package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"sparklassignment/backend/internal/dto"
	"sparklassignment/backend/internal/models"
)

type CMSService struct {
	db *gorm.DB
}

func NewCMSService(db *gorm.DB) *CMSService {
	return &CMSService{db: db}
}

func toQuizSummaryResponse(quiz models.Quiz) dto.QuizSummaryResponse {
	return dto.QuizSummaryResponse{
		QuizID:          quiz.QuizID,
		Category:        quiz.Category,
		Title:           quiz.Title,
		QuestionCount:   quiz.QuestionCount,
		TotalMarks:      quiz.TotalMarks,
		DurationMinutes: quiz.DurationMinutes,
		CreatedBy:       quiz.CreatedBy,
		CreatedAt:       quiz.CreatedAt,
		UpdatedAt:       quiz.UpdatedAt,
	}
}

func toQuestionResponse(question models.Question, options, correctOptions []string) dto.QuestionResponse {
	return dto.QuestionResponse{
		QuestionID:     question.QuestionID,
		Category:       question.Category,
		Prompt:         question.Prompt,
		Options:        options,
		CorrectOptions: correctOptions,
		Solution:       question.Solution,
	}
}

func (s *CMSService) GetQuizzes() (*dto.GetQuizzesResponse, error) {
	var quizzes []models.Quiz
	if err := s.db.Order("created_at desc").Find(&quizzes).Error; err != nil {
		return nil, err
	}

	response := make([]dto.QuizSummaryResponse, 0, len(quizzes))
	for _, quiz := range quizzes {
		response = append(response, toQuizSummaryResponse(quiz))
	}

	return &dto.GetQuizzesResponse{Quizzes: response}, nil
}

func (s *CMSService) GetQuestions() (*dto.GetQuestionsResponse, error) {
	var questions []models.Question
	if err := s.db.Order("question_id asc").Find(&questions).Error; err != nil {
		return nil, err
	}

	response := make([]dto.QuestionResponse, 0, len(questions))
	for _, question := range questions {
		options, err := ParseJSONStringArray(question.Options)
		if err != nil {
			return nil, fmt.Errorf("parse question options: %w", err)
		}

		correctOptions, err := ParseJSONStringArray(question.CorrectOptions)
		if err != nil {
			return nil, fmt.Errorf("parse correct options: %w", err)
		}

		response = append(response, toQuestionResponse(question, options, correctOptions))
	}

	return &dto.GetQuestionsResponse{Questions: response}, nil
}

func (s *CMSService) CreateQuestion(req dto.CreateQuestionRequest) (*dto.CreateQuestionResponse, error) {
	category := strings.TrimSpace(req.Category)
	prompt := strings.TrimSpace(req.Prompt)
	solution := strings.TrimSpace(req.Solution)

	if prompt == "" {
		return nil, ValidationError{Message: "prompt is required"}
	}
	if category == "" {
		return nil, ValidationError{Message: "category is required"}
	}
	if len(category) > 80 {
		return nil, ValidationError{Message: "category must be at most 80 characters"}
	}

	options := make([]string, 0, len(req.Options))
	seenOptions := make(map[string]struct{}, len(req.Options))
	for _, option := range req.Options {
		trimmed := strings.TrimSpace(option)
		if trimmed == "" {
			continue
		}

		if _, exists := seenOptions[trimmed]; exists {
			return nil, ValidationError{Message: "options must be unique"}
		}

		seenOptions[trimmed] = struct{}{}
		options = append(options, trimmed)
	}

	if len(options) < 2 {
		return nil, ValidationError{Message: "at least 2 options are required"}
	}

	correctOptions := make([]string, 0, len(req.CorrectOptions))
	seenCorrectOptions := make(map[string]struct{}, len(req.CorrectOptions))
	for _, correctOption := range req.CorrectOptions {
		trimmed := strings.TrimSpace(correctOption)
		if trimmed == "" {
			continue
		}

		if _, exists := seenOptions[trimmed]; !exists {
			return nil, ValidationError{Message: "correct_options must match the provided options"}
		}

		if _, exists := seenCorrectOptions[trimmed]; exists {
			return nil, ValidationError{Message: "correct_options must be unique"}
		}

		seenCorrectOptions[trimmed] = struct{}{}
		correctOptions = append(correctOptions, trimmed)
	}

	if len(correctOptions) == 0 {
		return nil, ValidationError{Message: "at least 1 correct option is required"}
	}

	rawOptions, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("marshal options: %w", err)
	}

	rawCorrectOptions, err := json.Marshal(correctOptions)
	if err != nil {
		return nil, fmt.Errorf("marshal correct options: %w", err)
	}

	question := models.Question{
		Category:       category,
		Prompt:         prompt,
		Options:        datatypes.JSON(rawOptions),
		CorrectOptions: datatypes.JSON(rawCorrectOptions),
		Solution:       solution,
	}

	if err := s.db.Create(&question).Error; err != nil {
		return nil, err
	}

	return &dto.CreateQuestionResponse{
		Message:  "question created successfully",
		Question: toQuestionResponse(question, options, correctOptions),
	}, nil
}

func (s *CMSService) GetQuiz(quizID uint) (*dto.GetQuizResponse, error) {
	var quiz models.Quiz
	if err := s.db.First(&quiz, quizID).Error; err != nil {
		return nil, err
	}

	var quizQuestions []models.QuizQuestion
	if err := s.db.Where("quiz_id = ?", quiz.QuizID).Order("sequence_number asc").Find(&quizQuestions).Error; err != nil {
		return nil, err
	}

	questionIDs := make([]uint, 0, len(quizQuestions))
	for _, quizQuestion := range quizQuestions {
		questionIDs = append(questionIDs, quizQuestion.QuestionID)
	}

	questionMap := make(map[uint]models.Question, len(questionIDs))
	if len(questionIDs) > 0 {
		var questions []models.Question
		if err := s.db.Where("question_id IN ?", questionIDs).Find(&questions).Error; err != nil {
			return nil, err
		}

		for _, question := range questions {
			questionMap[question.QuestionID] = question
		}
	}

	responseQuestions := make([]dto.QuizQuestionDetailResponse, 0, len(quizQuestions))
	for _, quizQuestion := range quizQuestions {
		prompt := ""
		if question, exists := questionMap[quizQuestion.QuestionID]; exists {
			prompt = question.Prompt
		}

		responseQuestions = append(responseQuestions, dto.QuizQuestionDetailResponse{
			QuizQuestionID: quizQuestion.QuizQuestionID,
			QuestionID:     quizQuestion.QuestionID,
			SequenceNumber: quizQuestion.SequenceNumber,
			Marks:          quizQuestion.Marks,
			Prompt:         prompt,
		})
	}

	return &dto.GetQuizResponse{
		Quiz:      toQuizSummaryResponse(quiz),
		Questions: responseQuestions,
	}, nil
}

func (s *CMSService) CreateQuiz(req dto.CreateQuizRequest, userID uint) (*dto.CreateQuizResponse, error) {
	category := strings.TrimSpace(req.Category)
	if category == "" {
		return nil, ValidationError{Message: "category is required"}
	}
	if len(category) > 80 {
		return nil, ValidationError{Message: "category must be at most 80 characters"}
	}
	if req.QuestionCount <= 0 || req.TotalMarks <= 0 || req.DurationMinutes <= 0 {
		return nil, ValidationError{Message: "question_count, total_marks and duration_minutes must be greater than zero"}
	}

	if len(req.Questions) != req.QuestionCount {
		return nil, ValidationError{Message: "number of questions in payload does not match question_count"}
	}

	totalQuestionMarks := 0
	seenSequenceNumbers := make(map[int]bool, len(req.Questions))
	seenQuestionIDs := make(map[uint]bool, len(req.Questions))

	for _, question := range req.Questions {
		if question.SequenceNumber <= 0 || question.Marks <= 0 {
			return nil, ValidationError{Message: "sequence_number and marks must be greater than zero"}
		}

		if seenSequenceNumbers[question.SequenceNumber] {
			return nil, ValidationError{Message: "sequence_number already in use"}
		}
		seenSequenceNumbers[question.SequenceNumber] = true

		if seenQuestionIDs[question.QuestionID] {
			return nil, ValidationError{Message: "question_id already in use"}
		}
		seenQuestionIDs[question.QuestionID] = true

		totalQuestionMarks += question.Marks
	}

	if totalQuestionMarks != req.TotalMarks {
		return nil, ValidationError{Message: "total question marks does not match total_marks"}
	}

	questionIDs := make([]uint, 0, len(req.Questions))
	for _, question := range req.Questions {
		questionIDs = append(questionIDs, question.QuestionID)
	}

	var existingQuestionsCount int64
	if err := s.db.Model(&models.Question{}).Where("question_id IN ?", questionIDs).Count(&existingQuestionsCount).Error; err != nil {
		return nil, err
	}

	if existingQuestionsCount != int64(len(questionIDs)) {
		return nil, ValidationError{Message: "one or more question_id does not exist"}
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	quiz := models.Quiz{
		Category:        category,
		Title:           req.Title,
		QuestionCount:   req.QuestionCount,
		TotalMarks:      req.TotalMarks,
		DurationMinutes: req.DurationMinutes,
		CreatedBy:       userID,
	}

	if err := tx.Create(&quiz).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	quizQuestions := make([]models.QuizQuestion, 0, len(req.Questions))
	for _, question := range req.Questions {
		quizQuestions = append(quizQuestions, models.QuizQuestion{
			QuizID:         quiz.QuizID,
			QuestionID:     question.QuestionID,
			SequenceNumber: question.SequenceNumber,
			Marks:          question.Marks,
		})
	}

	if err := tx.Create(&quizQuestions).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.CreateQuizResponse{
		Message: "quiz created successfully",
		Quiz:    toQuizSummaryResponse(quiz),
	}, nil
}

func (s *CMSService) GetParticipants(quizID uint) (*dto.GetParticipantsResponse, error) {
	var quiz models.Quiz
	if err := s.db.First(&quiz, quizID).Error; err != nil {
		return nil, err
	}

	var students []models.User
	if err := s.db.Where("role = ?", "student").Order("user_id asc").Find(&students).Error; err != nil {
		return nil, err
	}

	var attempts []models.QuizAttempt
	if err := s.db.Where("quiz_id = ?", quiz.QuizID).Order("updated_at desc").Find(&attempts).Error; err != nil {
		return nil, err
	}

	for index := range attempts {
		if attemptExpired(attempts[index], quiz) {
			if err := finalizeAttempt(s.db, &attempts[index]); err != nil {
				return nil, err
			}
		}
	}

	latestAttemptByStudent := make(map[uint]models.QuizAttempt)
	for _, attempt := range attempts {
		if _, exists := latestAttemptByStudent[attempt.StudentID]; !exists {
			latestAttemptByStudent[attempt.StudentID] = attempt
		}
	}

	response := make([]dto.ParticipantResponse, 0, len(students))
	for _, student := range students {
		item := dto.ParticipantResponse{
			StudentID: student.UserID,
			Username:  student.Username,
			FullName:  student.FullName,
			Status:    "not_started",
			Score:     0,
		}

		if attempt, exists := latestAttemptByStudent[student.UserID]; exists {
			attemptID := attempt.QuizAttemptID
			item.AttemptID = &attemptID
			item.Status = attempt.Status
			item.Score = attempt.Score

			if attempt.Status == "in_progress" {
				bundle, err := loadAttemptBundle(s.db, attempt)
				if err != nil {
					return nil, err
				}

				report, err := buildAttemptReport(bundle)
				if err != nil {
					return nil, err
				}

				item.Score = report.Score
			}
		}

		response = append(response, item)
	}

	return &dto.GetParticipantsResponse{
		Quiz:         toQuizSummaryResponse(quiz),
		Participants: response,
	}, nil
}

func (s *CMSService) GetParticipantReport(quizID, studentID uint) (*dto.AttemptReportResponse, error) {
	var attempt models.QuizAttempt
	if err := s.db.Where("quiz_id = ? AND student_id = ?", quizID, studentID).
		Order("updated_at desc").
		First(&attempt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, err
	}

	bundle, err := loadAttemptBundle(s.db, attempt)
	if err != nil {
		return nil, err
	}

	if attemptExpired(bundle.Attempt, bundle.Quiz) {
		if err := finalizeAttempt(s.db, &attempt); err != nil {
			return nil, err
		}

		bundle, err = loadAttemptBundle(s.db, attempt)
		if err != nil {
			return nil, err
		}
	}

	return buildAttemptReport(bundle)
}
