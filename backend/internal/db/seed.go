package db

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"sparklassignment/backend/internal/models"
)

func Seed(database *gorm.DB) error {
	if err := seedUsers(database); err != nil {
		return err
	}

	if err := seedQuestions(database); err != nil {
		return err
	}

	if err := seedQuiz(database); err != nil {
		return err
	}

	return nil
}

func seedUsers(database *gorm.DB) error {
	var count int64
	if err := database.Model(&models.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("seedUsers: %w", err)
	}
	if count > 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("seedUsers: %w", err)
	}

	users := []models.User{
		{
			Username:       "cmsadmin",
			HashedPassword: string(hashedPassword),
			Role:           "cms_admin",
			FullName:       "CMS Admin",
			Email:          "cmsadmin@example.com",
		},
		{
			Username:       "student1",
			HashedPassword: string(hashedPassword),
			Role:           "student",
			FullName:       "Student One",
			Email:          "student1@example.com",
		},
		{
			Username:       "student2",
			HashedPassword: string(hashedPassword),
			Role:           "student",
			FullName:       "Student Two",
			Email:          "student2@example.com",
		},
	}

	if err := database.Create(&users).Error; err != nil {
		return fmt.Errorf("seedUsers: %w", err)
	}
	return nil
}

func seedQuestions(database *gorm.DB) error {
	var count int64

	if err := database.Model(&models.Question{}).Count(&count).Error; err != nil {
		return fmt.Errorf("seedQuestions: %w", err)
	}
	if count > 0 {
		return nil
	}

	makeJSON := func(v any) datatypes.JSON {
		b, _ := json.Marshal(v)
		return datatypes.JSON(b)
	}

	questions := []models.Question{
		{
			Category:       "Geography",
			Prompt:         "What is the capital of France?",
			Options:        makeJSON([]string{"Berlin", "Madrid", "Paris", "Rome"}),
			CorrectOptions: makeJSON([]string{"Paris"}),
			Solution:       "Paris is the capital of France.",
		},
		{
			Category:       "Frontend",
			Prompt:         "Which language is used with React?",
			Options:        makeJSON([]string{"Go", "JavaScript", "Rust", "C"}),
			CorrectOptions: makeJSON([]string{"JavaScript"}),
			Solution:       "React is primarily used with JavaScript and TypeScript.",
		},
		{
			Category:       "Database",
			Prompt:         "What does SQL stand for?",
			Options:        makeJSON([]string{"Structured Query Language", "Simple Query Logic", "System Query Language", "Sequential Query Language"}),
			CorrectOptions: makeJSON([]string{"Structured Query Language"}),
			Solution:       "SQL stands for Structured Query Language.",
		},
		{
			Category:       "Programming",
			Prompt:         "Which of the following are programming languages?",
			Options:        makeJSON([]string{"Go", "Python", "HTML", "Java"}),
			CorrectOptions: makeJSON([]string{"Go", "Python", "Java"}),
			Solution:       "HTML is markup, while Go, Python, and Java are programming languages.",
		},
		{
			Category:       "Mathematics",
			Prompt:         "What is 2 + 2?",
			Options:        makeJSON([]string{"3", "4", "5", "6"}),
			CorrectOptions: makeJSON([]string{"4"}),
			Solution:       "2 + 2 equals 4.",
		},
	}

	if err := database.Create(&questions).Error; err != nil {
		return fmt.Errorf("seedQuestions: %w", err)
	}

	return nil
}

func seedQuiz(database *gorm.DB) error {
	var count int64
	if err := database.Model(&models.Quiz{}).Count(&count).Error; err != nil {
		return fmt.Errorf("seedQuiz: %w", err)
	}
	if count > 0 {
		return nil
	}

	var cmsUser models.User
	if err := database.Where("role = ?", "cms_admin").First(&cmsUser).Error; err != nil {
		return fmt.Errorf("find cms user: %w", err)
	}

	var questions []models.Question
	if err := database.Order("question_id asc").Find(&questions).Error; err != nil {
		return fmt.Errorf("find questions: %w", err)
	}
	if len(questions) < 3 {
		return fmt.Errorf("not enough questions to seed quiz")
	}

	quiz := models.Quiz{
		Category:        "Foundations",
		Title:           "Sample Quiz",
		QuestionCount:   3,
		TotalMarks:      30,
		DurationMinutes: 15,
		CreatedBy:       cmsUser.UserID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := database.Create(&quiz).Error; err != nil {
		return fmt.Errorf("create quiz: %w", err)
	}

	quizQuestions := []models.QuizQuestion{
		{QuizID: quiz.QuizID, QuestionID: questions[0].QuestionID, SequenceNumber: 1, Marks: 10},
		{QuizID: quiz.QuizID, QuestionID: questions[1].QuestionID, SequenceNumber: 2, Marks: 10},
		{QuizID: quiz.QuizID, QuestionID: questions[2].QuestionID, SequenceNumber: 3, Marks: 10},
	}

	if err := database.Create(&quizQuestions).Error; err != nil {
		return fmt.Errorf("create quiz questions: %w", err)
	}

	return nil
}
