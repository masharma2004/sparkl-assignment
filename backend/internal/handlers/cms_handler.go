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

type CMSHandler struct {
	service *services.CMSService
}

func NewCMSHandler(db *gorm.DB, _ *config.Config) *CMSHandler {
	return &CMSHandler{
		service: services.NewCMSService(db),
	}
}

func (h *CMSHandler) GetQuizzes(ctx *gin.Context) {
	response, err := h.service.GetQuizzes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch quizzes"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *CMSHandler) GetQuestions(ctx *gin.Context) {
	response, err := h.service.GetQuestions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch questions"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *CMSHandler) CreateQuestion(ctx *gin.Context) {
	var req dto.CreateQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid question payload"})
		return
	}

	response, err := h.service.CreateQuestion(req)
	if err != nil {
		var validationErr services.ValidationError
		if errors.As(err, &validationErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create question"})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *CMSHandler) GetQuiz(ctx *gin.Context) {
	quizID, ok := parseUintParam(ctx, "quiz_id")
	if !ok {
		return
	}

	response, err := h.service.GetQuiz(quizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch quiz"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *CMSHandler) CreateQuiz(ctx *gin.Context) {
	var req dto.CreateQuizRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz payload"})
		return
	}

	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	response, err := h.service.CreateQuiz(req, userID)
	if err != nil {
		var validationErr services.ValidationError
		if errors.As(err, &validationErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Message})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create quiz"})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *CMSHandler) GetParticipants(ctx *gin.Context) {
	quizID, ok := parseUintParam(ctx, "quiz_id")
	if !ok {
		return
	}

	response, err := h.service.GetParticipants(quizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch participants"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *CMSHandler) GetParticipantReport(ctx *gin.Context) {
	quizID, ok := parseUintParam(ctx, "quiz_id")
	if !ok {
		return
	}

	studentID, ok := parseUintParam(ctx, "student_id")
	if !ok {
		return
	}

	response, err := h.service.GetParticipantReport(quizID, studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "participant attempt not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch participant report"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
