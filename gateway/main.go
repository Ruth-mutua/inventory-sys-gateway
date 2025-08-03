package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-inventory-system/gateway/config"
	"go-inventory-system/gateway/middleware"
	"go-inventory-system/gateway/router"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize router with routes
	routes, err := config.LoadRoutes("routes.yaml")
	if err != nil {
		log.Fatalf("Failed to load routes: %v", err)
	}

	// Create router
	router := router.NewRouter(routes)

	// Setup middleware
	router.Use(middleware.Logging)
	router.Use(middleware.RateLimiting)
	router.Use(middleware.Metrics)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Gateway starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gateway...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Gateway exited")
}
