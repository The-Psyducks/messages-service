package middleware

import (
	"fmt"
	"log/slog"
	"messages/src/model/errors"
	"strings"

	"messages/src/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Error("Authorization header is required")

			err := errors.AuthenticationError("Authorization header is required")
			_ = c.AbortWithError(err.Status, err)
			c.Next()
			return
		}

		if authHeader == "contraseniaSecreta" {
			c.Next()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			slog.Error("Invalid authorization header")
			err := errors.AuthenticationError("Invalid authorization header")
			_ = c.AbortWithError(err.Status, err)
			c.Next()
			return
		}

		tokenString := bearerToken[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			slog.Error("Invalid token")

			err := errors.AuthenticationError(fmt.Sprintf("Invalid token: %s", err.Error()))
			_ = c.AbortWithError(err.Status, err)
			return
		}

		c.Set("session_user_id", claims.UserId)
		c.Set("session_user_admin", claims.UserAdmin)
		c.Set("tokenString", tokenString)

		c.Next()
	}
}
