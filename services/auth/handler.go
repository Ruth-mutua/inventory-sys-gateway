package main

import (
	"encoding/json"
	"net/http"

	"go-inventory-system/shared"

	"gorm.io/gorm"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req shared.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" || req.Username == "" {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Email, password, and username are required")
		return
	}

	// Check if user already exists
	var existingUser shared.User
	if err := h.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		shared.WriteErrorResponse(w, http.StatusConflict, "User already exists")
		return
	}

	// Hash password
	hashedPassword, err := shared.HashPassword(req.Password)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Create user
	user := shared.User{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := h.db.Create(&user).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := shared.GenerateJWT(user.ID, user.Email)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return response
	response := shared.AuthResponse{
		Token: token,
		User:  user,
	}

	shared.WriteSuccessResponse(w, http.StatusCreated, "User registered successfully", response)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req shared.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Find user
	var user shared.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	if !shared.CheckPassword(req.Password, user.Password) {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := shared.GenerateJWT(user.ID, user.Email)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return response
	response := shared.AuthResponse{
		Token: token,
		User:  user,
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Login successful", response)
}
