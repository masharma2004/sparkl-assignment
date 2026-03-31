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

type AuthHandler struct {
	service          *services.AuthService
	authCookieConfig config.AuthCookieConfig
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service:          services.NewAuthService(db, cfg),
		authCookieConfig: cfg.AuthCookieConfig(),
	}
}

const accessCookiePath = "/api/v1"
const refreshCookiePath = "/api/v1/auth"

func (h *AuthHandler) setAccessCookie(ctx *gin.Context, token string) {
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(
		h.authCookieConfig.AccessCookieName,
		token,
		h.authCookieConfig.AccessMaxAgeSeconds,
		accessCookiePath,
		h.authCookieConfig.Domain,
		h.authCookieConfig.Secure,
		true,
	)
}

func (h *AuthHandler) setRefreshCookie(ctx *gin.Context, token string) {
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(
		h.authCookieConfig.RefreshCookieName,
		token,
		h.authCookieConfig.RefreshMaxAgeSeconds,
		refreshCookiePath,
		h.authCookieConfig.Domain,
		h.authCookieConfig.Secure,
		true,
	)
}

func (h *AuthHandler) clearAccessCookie(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(
		h.authCookieConfig.AccessCookieName,
		"",
		-1,
		accessCookiePath,
		h.authCookieConfig.Domain,
		h.authCookieConfig.Secure,
		true,
	)
}

func (h *AuthHandler) clearRefreshCookie(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(
		h.authCookieConfig.RefreshCookieName,
		"",
		-1,
		refreshCookiePath,
		h.authCookieConfig.Domain,
		h.authCookieConfig.Secure,
		true,
	)
}

func (h *AuthHandler) setSessionCookies(ctx *gin.Context, session *services.SessionBundle) {
	h.setAccessCookie(ctx, session.AccessToken)
	h.setRefreshCookie(ctx, session.RefreshToken)
}

func (h *AuthHandler) clearSessionCookies(ctx *gin.Context) {
	h.clearAccessCookie(ctx)
	h.clearRefreshCookie(ctx)
}

func (h *AuthHandler) CMSLogin(ctx *gin.Context) {
	h.loginByRole(ctx, "cms_admin")
}

func (h *AuthHandler) StudentLogin(ctx *gin.Context) {
	h.loginByRole(ctx, "student")
}

func (h *AuthHandler) StudentSignup(ctx *gin.Context) {
	var req dto.StudentSignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username, password, full_name and email are required"})
		return
	}

	session, err := h.service.SignupStudent(req)
	if err != nil {
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

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student account"})
		return
	}

	h.setSessionCookies(ctx, session)
	ctx.JSON(http.StatusCreated, session.Response)
}

func (h *AuthHandler) loginByRole(ctx *gin.Context, role string) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	session, err := h.service.LoginByRole(req.Username, req.Password, role)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username/password"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	h.setSessionCookies(ctx, session)
	ctx.JSON(http.StatusOK, session.Response)
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(h.authCookieConfig.RefreshCookieName)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	session, err := h.service.RefreshSession(refreshToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRefreshSession) {
			h.clearSessionCookies(ctx)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh session is invalid"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh session"})
		return
	}

	h.setSessionCookies(ctx, session)
	ctx.JSON(http.StatusOK, session.Response)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	refreshToken, _ := ctx.Cookie(h.authCookieConfig.RefreshCookieName)
	if err := h.service.Logout(refreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	h.clearSessionCookies(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) Me(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user token"})
		return
	}

	response, err := h.service.GetMe(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
