package middleware

import (
	"exp_tracker/config"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ValidateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: missing or invalid token",
			})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		var userID int64
		var tokenExpiry time.Time
		query := `SELECT id, token_expires_at FROM users WHERE token = $1`

		err := config.Db.QueryRow(query, token).Scan(&userID, &tokenExpiry)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: invalid token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)

		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) (int64, error) {
	uid, exists := ctx.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("user_id not found in context")
	}
	userID, ok := uid.(int64)
	if !ok {
		return 0, fmt.Errorf("user_id has wrong type")
	}
	return userID, nil
}
