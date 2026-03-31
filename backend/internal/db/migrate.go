package db

import (
	"sparklassignment/backend/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(database *gorm.DB) error {
	if err := database.AutoMigrate(
		&models.User{},
		&models.Question{},
		&models.Quiz{},
		&models.QuizQuestion{},
		&models.QuizAttempt{},
		&models.AttemptAnswer{},
		&models.RefreshSession{},
		&models.RateLimitEntry{},
	); err != nil {
		return err
	}

	indexStatements := []string{
		`UPDATE questions SET category = 'General' WHERE category IS NULL OR trim(category) = ''`,
		`UPDATE quizzes SET category = 'General' WHERE category IS NULL OR trim(category) = ''`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_quiz_attempts_student_quiz_in_progress
			ON quiz_attempts (quiz_id, student_id)
			WHERE status = 'in_progress'`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_attempt_answers_attempt_question
			ON attempt_answers (attempt_id, quiz_question_id)`,
	}

	for _, statement := range indexStatements {
		if err := database.Exec(statement).Error; err != nil {
			return err
		}
	}

	return nil
}
