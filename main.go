package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/eventstore"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/transaction"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/router"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase"
)

func main() {
	fmt.Println("Starting Event Sourcing Todo Application...")

	// DI Container setup
	todoUseCase := usecase.NewTodoUseCase(tx, eventStore)
	todoHandler := handler.NewTodoHandler(todoUseCase)
	appRouter := router.NewRouter(todoHandler)

	// Setup routes
	mux := appRouter.SetupRoutes()

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	
	log.Fatal(http.ListenAndServe(port, mux))
}
