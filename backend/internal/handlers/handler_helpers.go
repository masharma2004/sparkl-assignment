package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getUserIDFromContext(ctx *gin.Context) (uint, bool) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user token"})
		return 0, false
	}

	return userID, true
}

func parseUintParam(ctx *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(ctx.Param(name), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid %s", name)})
		return 0, false
	}

	return uint(value), true
}
