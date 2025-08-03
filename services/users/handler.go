package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-inventory-system/shared"

	"gorm.io/gorm"
)

// UserHandler handles user-related requests
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// HandleUsers handles /users endpoint (GET, POST)
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleUser handles /users/{id} endpoint (GET, PUT, DELETE)
func (h *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userID, err := strconv.ParseUint(pathParts[2], 10, 32)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUser(w, r, uint(userID))
	case http.MethodPut:
		h.UpdateUser(w, r, uint(userID))
	case http.MethodDelete:
		h.DeleteUser(w, r, uint(userID))
	default:
		shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// ListUsers returns all users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []shared.User
	if err := h.db.Find(&users).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Users retrieved successfully", users)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user shared.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusCreated, "User created successfully", user)
}

// GetUser returns a specific user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, userID uint) {
	var user shared.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			shared.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		} else {
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "User retrieved successfully", user)
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, userID uint) {
	var user shared.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			shared.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		} else {
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		}
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.db.Model(&user).Updates(updateData).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "User updated successfully", user)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, userID uint) {
	var user shared.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			shared.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		} else {
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		}
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}

// GetCurrentUser returns the current authenticated user's profile
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var user shared.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			shared.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		} else {
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Profile retrieved successfully", user)
}
