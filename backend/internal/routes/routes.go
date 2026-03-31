package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sparklassignment/backend/internal/config"
	"sparklassignment/backend/internal/handlers"
	"sparklassignment/backend/internal/middleware"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "server is healthy"})
	})

	authHandler := handlers.NewAuthHandler(db, cfg)
	cmsHandler := handlers.NewCMSHandler(db, cfg)
	studentHandler := handlers.NewStudentHandler(db, cfg)
	authRateLimiter := middleware.NewSharedRateLimiter(db, 10, 5*time.Minute)
	refreshRateLimiter := middleware.NewSharedRateLimiter(db, 20, 5*time.Minute)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/cms/login", authRateLimiter.Middleware(), authHandler.CMSLogin)
			auth.POST("/student/login", authRateLimiter.Middleware(), authHandler.StudentLogin)
			auth.POST("/student/signup", authRateLimiter.Middleware(), authHandler.StudentSignup)
			auth.POST("/refresh", refreshRateLimiter.Middleware(), authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)

			auth.GET("/me",
				middleware.AuthMiddleware(cfg),
				authHandler.Me,
			)
		}

		cms := api.Group("/cms")
		cms.Use(middleware.AuthMiddleware(cfg), middleware.RequireRole("cms_admin"))
		{
			cms.GET("/questions", cmsHandler.GetQuestions)
			cms.POST("/questions", cmsHandler.CreateQuestion)
			cms.GET("/quizzes", cmsHandler.GetQuizzes)
			cms.GET("/quizzes/:quiz_id", cmsHandler.GetQuiz)
			cms.POST("/quizzes", cmsHandler.CreateQuiz)
			cms.GET("/quizzes/:quiz_id/participants", cmsHandler.GetParticipants)
			cms.GET("/quizzes/:quiz_id/participants/:student_id/report", cmsHandler.GetParticipantReport)
		}

		student := api.Group("/student")
		student.Use(middleware.AuthMiddleware(cfg), middleware.RequireRole("student"))
		{
			student.GET("/quizzes", studentHandler.GetQuizzes)
			student.GET("/attempts/:attempt_id", studentHandler.GetAttempt)
			student.GET("/attempts/:attempt_id/report", studentHandler.GetReport)
			student.PATCH("/attempts/:attempt_id/answers", studentHandler.SaveAnswer)
			student.POST("/attempts/:attempt_id/finish", studentHandler.FinishAttempt)

			student.POST("/quizzes/:quiz_id/start", studentHandler.StartQuiz)
		}
	}
}
