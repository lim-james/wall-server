package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    // "your_project/routers" // Import your project's router setup
)

// TestLogin checks the login process
func TestLogin(t *testing.T) {
    // Create a new HTTP request to your login endpoint
    user := map[string]string{
        "username": "user1",
        "password_hash": "password1hash",
    }
    requestBody, _ := json.Marshal(user)
    request, _ := http.NewRequest("POST", "/api/u/login", bytes.NewBuffer(requestBody))

    response := httptest.NewRecorder()
    handler := http.HandlerFunc(AuthHandler.LoginHandler) // Replace with your actual handler

    handler.ServeHTTP(response, request)

    if status := response.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // You might also want to check the response body for a JWT token or a success message
}