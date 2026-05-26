package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ozMadeBack/internal/models"

	"github.com/gin-gonic/gin"
)

func TestSellerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		user           any
		expectedStatus int
	}{
		{
			name:           "Unauthenticated",
			user:           nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Not a seller",
			user: models.User{
				IsSeller: false,
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Is a seller",
			user: models.User{
				IsSeller: true,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid user type",
			user:           "not a user struct",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.user != nil {
				c.Set("user", tt.user)
			}

			c.Request, _ = http.NewRequest("GET", "/", nil)

			handler := SellerMiddleware()
			handler(c)

			if w.Code != tt.expectedStatus && !(w.Code == 200 && tt.expectedStatus == 200) {
				// Note: gin.Context.Next() doesn't change status code by itself if no response is written
				// but AbortWithStatusJSON does.
			}

			if tt.expectedStatus != http.StatusOK {
				if w.Code != tt.expectedStatus {
					t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
				}
			} else {
				if c.IsAborted() {
					t.Errorf("expected request not to be aborted for valid seller")
				}
			}
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		role           any
		expectedStatus int
	}{
		{
			name:           "No role",
			role:           nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Not an admin",
			role:           "user",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Is an admin",
			role:           "admin",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid role type",
			role:           123,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.role != nil {
				c.Set("userRole", tt.role)
			}

			c.Request, _ = http.NewRequest("GET", "/", nil)

			handler := AdminMiddleware()
			handler(c)

			if tt.expectedStatus != http.StatusOK {
				if w.Code != tt.expectedStatus {
					t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
				}
			} else {
				if c.IsAborted() {
					t.Errorf("expected request not to be aborted for admin")
				}
			}
		})
	}
}
