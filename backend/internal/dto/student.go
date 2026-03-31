package dto

import "time"

type StudentQuizItemResponse struct {
	ID              uint   `json:"id"`
	Category        string `json:"category"`
	Title           string `json:"title"`
	QuestionCount   int    `json:"question_count"`
	TotalMarks      int    `json:"total_marks"`
	DurationMinutes int    `json:"duration_minutes"`
	Status          string `json:"status"`
	Action          string `json:"action"`
	AttemptID       *uint  `json:"attempt_id,omitempty"`
}

type StudentQuizResponse struct {
	Quizzes []StudentQuizItemResponse `json:"quizzes"`
}

type StartQuizResponse struct {
	Message   string `json:"message"`
	AttemptID uint   `json:"attempt_id"`
	Status    string `json:"status"`
}

type AttemptQuestionResponse struct {
	QuizQuestionID uint     `json:"quiz_question_id"`
	QuestionID     uint     `json:"question_id"`
	SequenceNumber int      `json:"sequence_number"`
	Marks          int      `json:"marks"`
	Prompt         string   `json:"prompt"`
	Options        []string `json:"options"`
	ChosenOptions  []string `json:"chosen_options"`
}

type AttemptQuizResponse struct {
	ID              uint   `json:"id"`
	Category        string `json:"category"`
	Title           string `json:"title"`
	DurationMinutes int    `json:"duration_minutes"`
	QuestionCount   int    `json:"question_count"`
	TotalMarks      int    `json:"total_marks"`
}

type GetAttemptResponse struct {
	AttemptID        uint                      `json:"attempt_id"`
	Status           string                    `json:"status"`
	RemainingSeconds int64                     `json:"remaining_seconds"`
	Quiz             AttemptQuizResponse       `json:"quiz"`
	Questions        []AttemptQuestionResponse `json:"questions"`
}

type SaveAnswerRequest struct {
	QuizQuestionID uint     `json:"quiz_question_id" binding:"required"`
	ChosenOptions  []string `json:"chosen_options"`
}

type SaveAnswerResponse struct {
	Message        string `json:"message"`
	AttemptID      uint   `json:"attempt_id"`
	QuizQuestionID uint   `json:"quiz_question_id"`
}

type FinishAttemptResponse struct {
	Message     string     `json:"message"`
	AttemptID   uint       `json:"attempt_id"`
	Status      string     `json:"status"`
	Score       int        `json:"score"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
}

type AttemptReportQuestionResponse struct {
	QuizQuestionID uint     `json:"quiz_question_id"`
	QuestionID     uint     `json:"question_id"`
	SequenceNumber int      `json:"sequence_number"`
	Prompt         string   `json:"prompt"`
	Options        []string `json:"options"`
	ChosenOptions  []string `json:"chosen_options"`
	CorrectOptions []string `json:"correct_options"`
	Solution       string   `json:"solution"`
	AwardedMarks   int      `json:"awarded_marks"`
	MaxMarks       int      `json:"max_marks"`
	IsCorrect      bool     `json:"is_correct"`
}

type AttemptReportResponse struct {
	AttemptID   uint                            `json:"attempt_id"`
	Status      string                          `json:"status"`
	Score       int                             `json:"score"`
	SubmittedAt *time.Time                      `json:"submitted_at,omitempty"`
	Quiz        AttemptQuizResponse             `json:"quiz"`
	Questions   []AttemptReportQuestionResponse `json:"questions"`
}
