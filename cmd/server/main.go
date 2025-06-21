package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mo-mohamed/acronis-memory-store/internal/api"
	"github.com/mo-mohamed/acronis-memory-store/internal/store"
	"github.com/mo-mohamed/acronis-memory-store/internal/store/memory"
)

func main() {
	// Create IStore instance
	var memoryStore store.IStore = memory.NewMemoryStore()

	// Create API handler
	handler := api.NewHandler(memoryStore)
	// Setup routes
	routes := handler.SetupRoutes()

	// Create HTTP server
	port := getEnvOrDefault("PORT", "8080")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: routes,
	}
	go func() {
		log.Printf("starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exited gracefully")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
