package middleware

import (
	"errors"
	"net/http"
	"sparklassignment/backend/internal/config"
	"sparklassignment/backend/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		tokenString := ""

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid Authorization header format",
				})
				return
			}

			tokenString = parts[1]
		} else {
			cookieValue, err := ctx.Cookie(cfg.AccessCookieName)
			if err != nil || strings.TrimSpace(cookieValue) == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Authentication required",
				})
				return
			}

			tokenString = cookieValue
		}

		token, err := jwt.ParseWithClaims(tokenString, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		claims, ok := token.Claims.(*utils.CustomClaims)
		if !ok || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		if claims.Issuer != cfg.JWTIssuer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token issuer",
			})
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
