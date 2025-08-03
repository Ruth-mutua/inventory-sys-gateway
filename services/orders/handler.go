package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-inventory-system/shared"

	"gorm.io/gorm"
)

// OrderHandler handles order-related requests
type OrderHandler struct {
	db *gorm.DB
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

// HandleOrders handles /orders endpoint (GET, POST)
func (h *OrderHandler) HandleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListOrders(w, r)
	case http.MethodPost:
		h.CreateOrder(w, r)
	default:
		shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleOrder handles /orders/{id} endpoint (GET, PUT, DELETE)
func (h *OrderHandler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	// Extract order ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	orderID, err := strconv.ParseUint(pathParts[2], 10, 32)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetOrder(w, r, uint(orderID))
	case http.MethodPut:
		h.UpdateOrder(w, r, uint(orderID))
	case http.MethodDelete:
		h.DeleteOrder(w, r, uint(orderID))
	default:
		shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetUserOrders handles /orders/user/{user_id} endpoint
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		shared.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract user ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userID, err := strconv.ParseUint(pathParts[3], 10, 32)
	if err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var orders []shared.Order
	if err := h.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "User orders retrieved successfully", orders)
}

// ListOrders returns all orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	var orders []shared.Order
	if err := h.db.Find(&orders).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Orders retrieved successfully", orders)
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order shared.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		shared.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Set user ID from context
	order.UserID = userID

	if err := h.db.Create(&order).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusCreated, "Order created successfully", order)
}

// GetOrder returns a specific order
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request, orderID uint) {
	var order shared.Order
	if err := h.db.First(&order, orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			shared.WriteErrorResponse(w, http.StatusNotFound, "Order not found")
		} else {
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch order")
		}
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Order retrieved successfully", order)
}

// UpdateOrder updates an order
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request, orderID uint) {
	var order shared.Order
	if err := h.db.First(&order, orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			shared.WriteErrorResponse(w, http.StatusNotFound, "Order not found")
		} else {
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch order")
		}
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		shared.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.db.Model(&order).Updates(updateData).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update order")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Order updated successfully", order)
}

// DeleteOrder deletes an order
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request, orderID uint) {
	var order shared.Order
	if err := h.db.First(&order, orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			shared.WriteErrorResponse(w, http.StatusNotFound, "Order not found")
		} else {
			shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch order")
		}
		return
	}

	if err := h.db.Delete(&order).Error; err != nil {
		shared.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete order")
		return
	}

	shared.WriteSuccessResponse(w, http.StatusOK, "Order deleted successfully", nil)
}
