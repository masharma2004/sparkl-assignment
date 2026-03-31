package middleware

import (
	"net/http"
	"sparklassignment/backend/internal/config"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := make(map[string]struct{}, len(cfg.AllowedOrigins))
	for _, origin := range cfg.AllowedOrigins {
		allowedOrigins[origin] = struct{}{}
	}

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		_, allowed := allowedOrigins[origin]

		if origin != "" && allowed {
			headers := ctx.Writer.Header()
			headers.Set("Access-Control-Allow-Origin", origin)
			headers.Set("Vary", "Origin")
			headers.Set("Access-Control-Allow-Credentials", "true")
			headers.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			headers.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			headers.Set("Access-Control-Max-Age", "600")
		}

		if ctx.Request.Method == http.MethodOptions {
			if origin != "" && !allowed {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Origin not allowed",
				})
				return
			}

			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
