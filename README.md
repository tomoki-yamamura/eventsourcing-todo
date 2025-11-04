# Event Sourcing Todo Application

This is a Todo application built using **Event Sourcing** architecture in Go.
It demonstrates Domain-Driven Design (DDD) principles, implementing CQRS pattern for command and query separation.

---

## Tech Stack

This project uses the following core technologies:

- [golang](https://go.dev/)
- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router
- [goose](https://github.com/pressly/goose) - Database migration tool
- [MySQL](https://www.mysql.com/) - Event store database
- [Docker](https://www.docker.com/) - Containerization

---

## Architecture

The application follows **Clean Architecture** principles with the following layers:

- **Domain Layer**: Aggregates, Value Objects, Commands, Events, and Repository interfaces
- **UseCase Layer**: Business logic, application services, and port interfaces
- **Infrastructure Layer**: Database implementations, HTTP handlers, presenters, views, and external integrations

### Clean Architecture Pattern

- **Ports & Adapters**: Clear separation between business logic and external concerns
- **Presenter Pattern**: UseCase outputs to Presenter interfaces, maintaining dependency inversion
- **View Layer**: HTTP response rendering separated from presentation logic

### Event Sourcing Components

- **Aggregates**: TodoListAggregate manages todo list state through events
- **Events**: TodoListCreatedEvent, TodoAddedEvent capture state changes
- **Event Store**: Persists events with optimistic locking for concurrency control
- **Read Models**: Separate query models for retrieving todo lists

---

### Key Features

- **Value Objects**: UserID and TodoText with built-in validation
- **Clean Architecture**: Strict separation of domain, use case, and infrastructure
- **Event Sourcing**: All state changes captured as immutable events
- **CQRS**: Command and query responsibility segregation
- **Domain Rules**: Business logic like "max 3 todos per day" enforced in domain layer
- **Optimistic Locking**: Prevents concurrent modification conflicts
- **Transaction Management**: Flexible transaction control with retry logic

---

## Requirements

- `direnv`
- `docker` and `docker-compose`
- `go` 1.21+

---

## Environment Setup

Install task:

```bash
brew install go-task/tap/go-task
```

or

```bash
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
```

Copy the example env file:

```bash
cp .envrc.example .envrc
```

Set up direnv:

```bash
direnv allow
```

Start MySQL database:

```bash
task docker:up
```

Wait for MySQL to be ready, then run database migrations:

```bash
task migrate:up
```

Create test database and grant permissions:

```bash
mysql -h 127.0.0.1 -P 23306 -u root -p${MYSQL_ROOT_PASSWORD} \
  -e "CREATE DATABASE IF NOT EXISTS event_test; GRANT ALL PRIVILEGES ON event_test.* TO 'test'@'%'; FLUSH PRIVILEGES;"
```

Setup test database with migrations:

```bash
task migrate:test:up
```

Verify everything is working by running tests:

```bash
task test
```

---

## API Endpoints

The application exposes RESTful APIs for managing todo lists:

### Create Todo List

```bash
POST /todo-lists
```

Request body:

```json
{
  "user_id": "user123"
}
```

### Add Todo Item

```bash
POST /todo-lists/{aggregate_id}/items
```

Request body:

```json
{
  "user_id": "user123",
  "text": "Learn Event Sourcing"
}
```

### Get Todo List

```bash
GET /todo-lists/{aggregate_id}/items
```

---

## Run Application

Start the API server:

```bash
task run
```

Once running, test the API:

1. Create a new todo list:

```bash
curl -X POST "http://localhost:8080/todo-lists" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123"}'
```

2. Add a todo item (replace {aggregate_id} with the ID from step 1):

```bash
curl -X POST "http://localhost:8080/todo-lists/{aggregate_id}/items" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123", "text": "Learn Event Sourcing"}'
```

3. Get todo list:

```bash
curl -X GET "http://localhost:8080/todo-lists/{aggregate_id}/items"
```

---

## Testing

Run tests:

```bash
task test
```

---

## Development

### Code Quality

Run linting:

```bash
task lint
```

---

## Project Structure

```
.
├── main.go                              # Application entry point
├── Taskfile.yaml                        # Task automation (build, test, migration)
├── docker-compose.yaml                  # MySQL database setup
├── .envrc.example                       # Environment configuration template
├── go.mod                               # Go module dependencies
├── container/                           # Dependency injection container
│   └── container.go                     # DI container implementation
├── docker/                              # Docker configuration files
└── internal/                            # Private application code
    ├── config/                          # Configuration management
    │   └── config.go                    # Application configuration
    ├── domain/                          # Domain layer (business logic)
    │   ├── aggregate/                   # Domain aggregates
    │   │   ├── todo_list.go             # TodoListAggregate implementation
    │   │   └── todo_list_test.go        # Aggregate unit tests
    │   ├── command/                     # Domain commands
    │   │   ├── create_todo_list_command.go
    │   │   └── add_todo_command.go
    │   ├── entity/                      # Domain entities
    │   │   └── todo_item.go             # TodoItem entity
    │   ├── event/                       # Domain events
    │   │   ├── event.go                 # Event interface
    │   │   ├── todo_list_created_event.go
    │   │   └── todo_added_event.go
    │   ├── repository/                  # Repository interfaces
    │   │   ├── event_store.go           # Event store interface
    │   │   ├── event_deserializer.go    # Event deserializer interface
    │   │   └── transaction.go           # Transaction interface
    │   └── value/                       # Value objects
    │       ├── user_id.go               # UserID value object with validation
    │       ├── user_id_test.go          # UserID tests
    │       ├── todo_text.go             # TodoText value object with validation
    │       ├── todo_text_test.go        # TodoText tests
    │       └── aggregate_id.go          # AggregateID value object
    ├── usecase/                         # Application layer (CQRS use cases)
    │   ├── command/                     # Command side (write operations)
    │   │   ├── input/                   # Command input models
    │   │   │   ├── create_todo_list_input.go
    │   │   │   └── add_todo_input.go
    │   │   ├── todo_list_create_command.go  # TodoList creation command
    │   │   └── todo_add_item_command.go     # TodoItem addition command
    │   ├── query/                       # Query side (read operations)
    │   │   ├── dto/                     # Data Transfer Objects
    │   │   │   └── todo_list_view_dto.go # TodoList view DTOs
    │   │   ├── input/                   # Query input models
    │   │   │   └── get_todo_list_input.go
    │   │   ├── output/                  # Query output models
    │   │   │   └── get_todo_list_output.go
    │   │   └── todo_list_query.go       # TodoList query implementation
    │   ├── ports/                       # Port interfaces (Clean Architecture)
    │   │   ├── gateway/                 # Gateway interfaces (external systems)
    │   │   │   ├── eventbus.go          # Event bus interfaces
    │   │   │   └── projector.go         # Projector interface
    │   │   ├── presenter/               # Presenter interfaces (output ports)
    │   │   │   └── todo_list_presenter.go # TodoList presenter interface
    │   │   └── readmodelstore/          # Read model store interfaces
    │   │       ├── todo_list_query.go   # Read model store interface
    │   │       └── dto/                 # Read model DTOs
    │   │           └── todo_list_view_dto.go
    │   └── readmodelstore/              # Read model store interface definitions
    │       └── todo_list_query.go       # Read model store interface
    └── infrastructure/                  # Infrastructure layer
        ├── bus/                         # Event bus implementation
        │   └── eventbus.go              # In-memory event bus
        ├── database/                    # Database implementations
        │   ├── client/                  # Database client implementation
        │   │   └── client.go            # MySQL client
        │   ├── eventstore/              # Event store implementation
        │   │   ├── event_store_impl.go  # EventStore implementation
        │   │   ├── event_store_impl_test.go # EventStore integration tests
        │   │   ├── migration/           # Database migration files
        │   │   └── deserializer/        # Event deserializers
        │   │       ├── event_deserializer_impl.go
        │   │       ├── todo_list_created_deserializer.go
        │   │       └── todo_added_deserializer.go
        │   ├── transaction/             # Transaction management
        │   │   └── transaction.go       # Transaction implementation
        │   └── testutil/                # Database test utilities
        │       └── setup_test_db.go     # Test database setup
        ├── projector/                   # Read model projectors
        │   └── todo/                    # Todo-specific projector
        │       ├── todo_projector.go    # Todo projector implementation
        │       ├── todo_projector_test.go # Projector unit tests
        │       ├── inmemory_repository.go # In-memory read model repository
        │       └── inmemory_repository_test.go # Repository unit tests
        ├── presenter/                   # Presenter layer (Clean Architecture)
        │   ├── viewmodel/               # View models for presentation
        │   │   └── todo_list_view_model.go # TodoList view model
        │   ├── view.go                  # View interface definition
        │   └── todo_presenter_impl.go   # TodoList presenter implementation
        ├── view/                        # View implementations
        │   └── todo_list_view_http.go   # HTTP view implementation
        ├── handler/                     # HTTP handlers (separated by responsibility)
        │   ├── command/                 # Command handlers (write operations)
        │   │   ├── todo_list_create_command_handler.go  # TodoList creation
        │   │   └── todo_add_item_command_handler.go     # TodoItem addition
        │   ├── query/                   # Query handlers (read operations)
        │   │   └── todo_list_query_handler.go           # TodoList queries
        │   ├── request/                 # HTTP request models
        │   │   └── todo_request.go
        │   └── response/                # HTTP response models
        │       └── todo_response.go
        └── router/                      # HTTP routing configuration
            └── router.go                # Router setup with separated handlers
```

---

## License

This project is for educational purposes demonstrating Event Sourcing patterns in Go.
