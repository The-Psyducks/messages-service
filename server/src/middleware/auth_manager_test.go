package middleware

import (
	"fmt"
	"log"
	"messages/src/auth"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fmt.Println("jwt_secret: ", os.Getenv("JWT_SECRET"))
	validBearerToken, err := auth.GenerateToken("userId", "userName", false)
	println("validBearerToken: ", validBearerToken)
	if err != nil {
		log.Fatalln("Error generating token: ", err)
	}
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "No Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header is required",
		},
		{
			name:           "Invalid Authorization Header",
			authHeader:     "InvalidHeader",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid authorization header",
		},
		{
			name:           "Invalid Token",
			authHeader:     "Bearer invalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
		{
			name:           "Valid Token",
			authHeader:     "Bearer " + validBearerToken,
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Secret Password",
			authHeader:     "contraseniaSecreta",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gin router
			r := gin.New()
			r.Use(AuthMiddleware())
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tt.authHeader)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			// Check the response
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}
		})
	}
}
