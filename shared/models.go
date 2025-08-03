package shared

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // Hidden from JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserClaims represents JWT claims for user authentication
type UserClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Order represents an order in the system
type Order struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	ProductName string    `json:"product_name" gorm:"not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	TotalPrice  float64   `json:"total_price" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:'pending'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AuthRequest represents login/register request
type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username,omitempty" validate:"omitempty,min=3"`
}

// AuthResponse represents login response
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Route represents a gateway route configuration
type Route struct {
	Path    string   `yaml:"path"`
	Backend string   `yaml:"backend"`
	Methods []string `yaml:"methods"`
}
