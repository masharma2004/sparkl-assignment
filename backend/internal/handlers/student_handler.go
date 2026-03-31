package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sparklassignment/backend/internal/config"
	"sparklassignment/backend/internal/dto"
	"sparklassignment/backend/internal/services"
)

type StudentHandler struct {
	service *services.StudentService
}

func NewStudentHandler(db *gorm.DB, _ *config.Config) *StudentHandler {
	return &StudentHandler{
		service: services.NewStudentService(db),
	}
}

func (h *StudentHandler) GetQuizzes(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	response, err := h.service.GetQuizzes(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch quizzes"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *StudentHandler) StartQuiz(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	quizID, ok := parseUintParam(ctx, "quiz_id")
	if !ok {
		return
	}

	response, state, err := h.service.StartQuiz(userID, quizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start quiz"})
		return
	}

	if state == "completed" {
		ctx.JSON(http.StatusConflict, gin.H{
			"error":      response.Message,
			"attempt_id": response.AttemptID,
		})
		return
	}

	if state == "resume" {
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *StudentHandler) GetAttempt(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	attemptID, ok := parseUintParam(ctx, "attempt_id")
	if !ok {
		return
	}

	response, err := h.service.GetAttempt(userID, attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "quiz attempt not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch quiz attempt"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *StudentHandler) SaveAnswer(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	attemptID, ok := parseUintParam(ctx, "attempt_id")
	if !ok {
		return
	}

	var req dto.SaveAnswerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid answer payload"})
		return
	}

	response, err := h.service.SaveAnswer(userID, attemptID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "quiz attempt not found"})
			return
		}

		var validationErr services.ValidationError
		if errors.As(err, &validationErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
			return
		}

		var conflictErr services.ConflictError
		if errors.As(err, &conflictErr) {
			ctx.JSON(http.StatusConflict, gin.H{"error": conflictErr.Message})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save answer"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *StudentHandler) FinishAttempt(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	attemptID, ok := parseUintParam(ctx, "attempt_id")
	if !ok {
		return
	}

	response, err := h.service.FinishAttempt(userID, attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "quiz attempt not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finish attempt"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *StudentHandler) GetReport(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	attemptID, ok := parseUintParam(ctx, "attempt_id")
	if !ok {
		return
	}

	response, err := h.service.GetReport(userID, attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "quiz attempt not found"})
			return
		}

		var conflictErr services.ConflictError
		if errors.As(err, &conflictErr) {
			ctx.JSON(http.StatusConflict, gin.H{"error": conflictErr.Message})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build attempt report"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
