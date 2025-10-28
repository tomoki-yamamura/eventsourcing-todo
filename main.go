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

	if err := cont.TodoProjector.Start(ctx, cont.EventBus); err != nil {
		log.Fatalf("Failed to start projector: %v", err)
	}

	// Handler layer setup (CQRS)
	createCommandHandler := command.NewTodoListCreateCommandHandler(cont.TodoListCreateCommand)
	addCommandHandler := command.NewTodoAddItemCommandHandler(cont.TodoAddItemCommand)
	queryHandler := query.NewTodoListQueryHandler(cont.QueryUseCase)

	// Router setup
	appRouter := router.NewRouter(createCommandHandler, addCommandHandler, queryHandler)
	mux := appRouter.SetupRoutes()

	// Start server
	port := ":" + cfg.HTTPPort
	fmt.Printf("Server starting on port %s\n", port)

	log.Fatal(http.ListenAndServe(port, mux))
}
