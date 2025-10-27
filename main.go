package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/tomoki-yamamura/eventsourcing-todo/container"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/config"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/query"
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

	// Handler layer setup (CQRS)
	commandHandler := command.NewTodoCommandHandler(cont.CommandUseCase)
	queryHandler := query.NewTodoQueryHandler(cont.QueryUseCase)
	
	// Router setup
	appRouter := router.NewRouter(commandHandler, queryHandler)
	mux := appRouter.SetupRoutes()

	// Start server
	port := ":" + cfg.HTTPPort
	fmt.Printf("Server starting on port %s\n", port)

	log.Fatal(http.ListenAndServe(port, mux))
}
