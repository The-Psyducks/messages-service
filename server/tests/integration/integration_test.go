package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"messages/src/auth"   // Import your auth package
	"messages/src/router" // Import your router package
)

func TestSendMessageHappyPath(t *testing.T) {
	// Set up the mock configuration
	gin.SetMode(gin.TestMode)

	// Create the router with mock configuration
	r, err := router.NewRouter(router.MOCK_EXTERNAL) // Assuming you have this method
	if err != nil {
		t.Fatalf("could not create router: %v", err)
	}

	userId := "userId"
	username := "testuser"
	token, err := auth.GenerateToken(userId, username, false)
	if err != nil {
		t.Fatalf("could not generate token: %v", err)
	}

	payload := map[string]string{
		"sender_id":   userId,          // Valid userId
		"receiver_id": userId,          // Valid userId
		"content":     "Hello, World!", // Valid content
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("could not marshal JSON: %v", err)
	}

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	// Check the response body
	var responseBody struct {
		ChatReference string `json:"chat-reference"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	// Verify the expected response
	expectedChatReference := "mockMessageRef" // Adjust based on the mock response in SendMessage
	if responseBody.ChatReference != expectedChatReference {
		t.Errorf("expected chat-reference to be %s, got %s", expectedChatReference, responseBody.ChatReference)
	}
}
