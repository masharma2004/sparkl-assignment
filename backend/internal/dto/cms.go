package dto

import "time"

type QuizSummaryResponse struct {
	QuizID          uint      `json:"quiz_id"`
	Category        string    `json:"category"`
	Title           string    `json:"title"`
	QuestionCount   int       `json:"question_count"`
	TotalMarks      int       `json:"total_marks"`
	DurationMinutes int       `json:"duration_minutes"`
	CreatedBy       uint      `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type GetQuizzesResponse struct {
	Quizzes []QuizSummaryResponse `json:"quizzes"`
}

type QuestionResponse struct {
	QuestionID     uint     `json:"question_id"`
	Category       string   `json:"category"`
	Prompt         string   `json:"prompt"`
	Options        []string `json:"options"`
	CorrectOptions []string `json:"correct_options,omitempty"`
	Solution       string   `json:"solution,omitempty"`
}

type GetQuestionsResponse struct {
	Questions []QuestionResponse `json:"questions"`
}

type CreateQuestionRequest struct {
	Category       string   `json:"category" binding:"required"`
	Prompt         string   `json:"prompt" binding:"required"`
	Options        []string `json:"options" binding:"required"`
	CorrectOptions []string `json:"correct_options" binding:"required"`
	Solution       string   `json:"solution"`
}

type CreateQuestionResponse struct {
	Message  string           `json:"message"`
	Question QuestionResponse `json:"question"`
}

type CreateQuizQuestionInput struct {
	QuestionID     uint `json:"question_id" binding:"required"`
	SequenceNumber int  `json:"sequence_number" binding:"required"`
	Marks          int  `json:"marks" binding:"required"`
}

type CreateQuizRequest struct {
	Category        string                    `json:"category" binding:"required"`
	Title           string                    `json:"title" binding:"required"`
	QuestionCount   int                       `json:"question_count" binding:"required"`
	TotalMarks      int                       `json:"total_marks" binding:"required"`
	DurationMinutes int                       `json:"duration_minutes" binding:"required"`
	Questions       []CreateQuizQuestionInput `json:"questions" binding:"required"`
}

type CreateQuizResponse struct {
	Message string              `json:"message"`
	Quiz    QuizSummaryResponse `json:"quiz"`
}

type QuizQuestionDetailResponse struct {
	QuizQuestionID uint   `json:"quiz_question_id"`
	QuestionID     uint   `json:"question_id"`
	SequenceNumber int    `json:"sequence_number"`
	Marks          int    `json:"marks"`
	Prompt         string `json:"prompt"`
}

type GetQuizResponse struct {
	Quiz      QuizSummaryResponse          `json:"quiz"`
	Questions []QuizQuestionDetailResponse `json:"questions"`
}

type ParticipantResponse struct {
	StudentID uint   `json:"student_id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Status    string `json:"status"`
	Score     int    `json:"score"`
	AttemptID *uint  `json:"attempt_id,omitempty"`
}

type GetParticipantsResponse struct {
	Quiz         QuizSummaryResponse   `json:"quiz"`
	Participants []ParticipantResponse `json:"participants"`
}
