package integration

//import (
//	"bytes"
//	"encoding/json"
//	"github.com/gin-gonic/gin"
//	"messages/src/auth"
//	"messages/src/router"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func AddDeviceTokenForUser(t *testing.T) {
//	// Set up the mock configuration
//	gin.SetMode(gin.TestMode)
//
//	// Create the router with mock configuration
//	r, err := router.NewRouter(router.MOCK_EXTERNAL) // Assuming you have this method
//	if err != nil {
//		t.Fatalf("could not create router: %v", err)
//	}
//
//	userId := "userId"
//	username := "testuser"
//	token, err := auth.GenerateToken(userId, username, false)
//	if err != nil {
//		t.Fatalf("could not generate token: %v", err)
//	}
//
//	payload := map[string]string{
//		"user_id":      userId, // Valid userId
//		"device_token": "deviceToken",
//	}
//
//	jsonPayload, err := json.Marshal(payload)
//	if err != nil {
//		t.Fatalf("could not marshal JSON: %v", err)
//	}
//
//
//	req, _ := http.NewRequest(http.MethodPost, "/device", bytes.NewBuffer(jsonPayload))
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+token)
//
//	// Create a response recorder
//	w := httptest.NewRecorder()
//
//	// Perform the request
//	r.ServeHTTP(w, req)
//
//
//	if w.Code != http.StatusOK {
//		t.Errorf("expected status code 200, got %d", w.Code)
//	}
//}
