package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SharedRateLimiter struct {
	db     *gorm.DB
	limit  int
	window time.Duration
}

type rateLimitResult struct {
	HitCount int
}

func NewSharedRateLimiter(db *gorm.DB, limit int, window time.Duration) *SharedRateLimiter {
	return &SharedRateLimiter{
		db:     db,
		limit:  limit,
		window: window,
	}
}

func (l *SharedRateLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		allowed, err := l.allow(ctx.ClientIP(), ctx.FullPath())
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limiter unavailable",
			})
			return
		}

		if allowed {
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many requests. Please try again later.",
		})
	}
}

func (l *SharedRateLimiter) allow(ipAddress, route string) (bool, error) {
	now := time.Now().UTC()
	resetAt := now.Add(l.window)
	routeKey := route
	if routeKey == "" {
		routeKey = "unknown"
	}

	scopeKey := fmt.Sprintf("%s|%s", ipAddress, routeKey)
	result := rateLimitResult{}

	err := l.db.Raw(`
		INSERT INTO rate_limit_entries (scope_key, hit_count, reset_at, created_at, updated_at)
		VALUES (?, 1, ?, ?, ?)
		ON CONFLICT (scope_key) DO UPDATE
		SET hit_count = CASE
				WHEN rate_limit_entries.reset_at <= ? THEN 1
				ELSE rate_limit_entries.hit_count + 1
			END,
			reset_at = CASE
				WHEN rate_limit_entries.reset_at <= ? THEN ?
				ELSE rate_limit_entries.reset_at
			END,
			updated_at = ?
		RETURNING hit_count
	`, scopeKey, resetAt, now, now, now, now, resetAt, now).Scan(&result).Error
	if err != nil {
		return false, err
	}

	return result.HitCount <= l.limit, nil
}
