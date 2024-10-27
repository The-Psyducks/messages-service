package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestSendMessageInvalidUser(t *testing.T) {
	// Set up the mock configuration
	gin.SetMode(gin.TestMode)

	// Create the router with mock configuration
	r, err := router.NewRouter(router.MOCK_EXTERNAL)
	if err != nil {
		t.Fatalf("could not create router: %v", err)
	}

	invalidUserId := "fakeUserId" // This ID should cause the mock CheckUserExists to return false
	username := "invaliduser"
	token, err := auth.GenerateToken(invalidUserId, username, false)
	if err != nil {
		t.Fatalf("could not generate token: %v", err)
	}

	payload := map[string]string{
		"sender_id":   invalidUserId,
		"receiver_id": "userId", // Valid userId to make sure sender is the only invalid part
		"content":     "Hello, World!",
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

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code 404, got %d", w.Code)
	}

	// Check the error message in the response body
	var responseBody struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Status int    `json:"status"`
		Detail string `json:"detail"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	// Verify the expected error message
	expectedErrorMessage := http.StatusBadRequest
	if responseBody.Status != expectedErrorMessage {
		t.Errorf("expected error message to be '%d', got '%d'", expectedErrorMessage, responseBody.Status)
	}
}

func TestSendMessageInvalidToken(t *testing.T) {
	// Set up the mock configuration
	gin.SetMode(gin.TestMode)

	// Create the router with mock configuration
	r, err := router.NewRouter(router.MOCK_EXTERNAL)
	if err != nil {
		t.Fatalf("could not create router: %v", err)
	}

	userId := "userId"
	username := "testuser"
	// Here, we'll generate a valid token but we will replace it with an invalid one.
	token, err := auth.GenerateToken(userId, username, false)
	if err != nil {
		t.Fatalf("could not generate token: %v", err)
	}

	// Intentionally modifying the token to make it invalid
	invalidToken := token + "invalid"

	payload := map[string]string{
		"sender_id":   userId,          // Valid userId
		"receiver_id": "userId",        // Valid userId
		"content":     "Hello, World!", // Valid content
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("could not marshal JSON: %v", err)
	}

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/messages", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+invalidToken) // Using the invalid token

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code 401, got %d", w.Code)
	}

	// Check the error message in the response body
	var responseBody struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Status int    `json:"status"`
		Detail string `json:"detail"`
	}
	fmt.Println("Response:", w.Body.String())
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	// Verify the expected error message
	expectedErrorTitle := "Authentication Error"
	if responseBody.Title != expectedErrorTitle {
		t.Errorf("expected error title to be '%s', got '%s'", expectedErrorTitle, responseBody.Title)
	}
	expectedErrorStatus := http.StatusUnauthorized
	if responseBody.Status != expectedErrorStatus {
		t.Errorf("expected error status to be '%d', got '%d'", expectedErrorStatus, responseBody.Status)
	}
}
