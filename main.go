package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/tomoki-yamamura/eventsourcing-todo/container"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/config"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/router"
)

func main() {
	fmt.Println("Starting Event Sourcing Todo Application")

	// Load config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// DI Container setup
	ctx := context.Background()
	cont := container.NewContainer()
	if err := cont.Inject(ctx, cfg); err != nil {
		log.Fatalf("Failed to inject dependencies: %v", err)
	}

	// Infrastructure layer setup
	todoHandler := handler.NewTodoHandler(cont.TodoUseCase)
	appRouter := router.NewRouter(todoHandler)

	// Setup routes
	mux := appRouter.SetupRoutes()

	// Start server
	port := ":" + cfg.HTTPPort
	fmt.Printf("Server starting on port %s\n", port)
	
	log.Fatal(http.ListenAndServe(port, mux))
}
