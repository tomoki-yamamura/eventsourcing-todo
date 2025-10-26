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
- **UseCase Layer**: Business logic and application services
- **Infrastructure Layer**: Database implementations, HTTP handlers, and external integrations

### Event Sourcing Components

- **Aggregates**: TodoListAggregate manages todo list state through events
- **Events**: TodoListCreatedEvent, TodoAddedEvent capture state changes
- **Event Store**: Persists events with optimistic locking for concurrency control
- **Read Models**: Separate query models for retrieving todo lists

---

## Project Structure

```
.
├── main.go                           # Application entry point
├── Taskfile.yaml                     # Task automation (build, test, migration)
├── docker-compose.yaml              # MySQL database setup
├── .envrc.example                   # Environment configuration template
├── go.mod                           # Go module dependencies
├── container/                       # Dependency injection container
├── docker/                         # Docker configuration files
└── internal/                       # Private application code
    ├── config/                     # Configuration management
    ├── domain/                     # Domain layer (business logic)
    │   ├── aggregate/              # Domain aggregates
    │   │   ├── todo_list.go        # TodoListAggregate implementation
    │   │   └── todo_list_test.go   # Aggregate unit tests
    │   ├── command/                # Domain commands
    │   │   ├── create_todo_list_command.go
    │   │   └── add_todo_command.go
    │   ├── entity/                 # Domain entities
    │   │   └── todo_item.go        # TodoItem entity
    │   ├── event/                  # Domain events
    │   │   ├── event.go            # Event interface
    │   │   ├── todo_list_created_event.go
    │   │   └── todo_added_event.go
    │   ├── repository/             # Repository interfaces
    │   │   ├── event_store.go      # Event store interface
    │   │   ├── transaction.go      # Transaction interface
    │   │   └── database_client.go  # Database client interface
    │   └── value/                  # Value objects
    │       ├── user_id.go          # UserID value object with validation
    │       ├── user_id_test.go     # UserID tests
    │       ├── todo_text.go        # TodoText value object with validation
    │       └── todo_text_test.go   # TodoText tests
    ├── usecase/                    # Application layer (use cases)
    │   ├── todo.go                 # Todo use case implementation
    │   ├── input/                  # Use case input models
    │   │   ├── create_todo_list_input.go
    │   │   ├── add_todo_input.go
    │   │   └── get_todo_list_input.go
    │   └── output/                 # Use case output models
    │       ├── create_todo_list_output.go
    │       └── get_todo_list_output.go
    └── infrastructure/             # Infrastructure layer
        ├── database/               # Database implementations
        │   ├── client/             # Database client implementation
        │   ├── eventstore/         # Event store implementation
        │   │   ├── event_store_impl.go      # EventStore implementation
        │   │   ├── event_store_impl_test.go # EventStore integration tests
        │   │   ├── migration/      # Database migration files
        │   │   └── deserializer/   # Event deserializers
        │   ├── transaction/        # Transaction management
        │   └── testutil/          # Database test utilities
        ├── handler/               # HTTP handlers
        │   ├── todo_handler.go    # Todo REST API handlers
        │   ├── request/           # HTTP request models
        │   └── response/          # HTTP response models
        └── router/                # HTTP routing configuration
```

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
POST /todo-lists/{aggregate_id}/todos
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
GET /todo-lists/{aggregate_id}/todos
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
curl -X POST "http://localhost:8080/todo-lists/{aggregate_id}/todos" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123", "text": "Learn Event Sourcing"}'
```

3. Get todo list:
```bash
curl -X GET "http://localhost:8080/todo-lists/{aggregate_id}/todos"
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
├── cmd/                    # Application entry points
├── internal/
│   ├── domain/            # Domain layer (aggregates, events, value objects)
│   │   ├── aggregate/     # Domain aggregates
│   │   ├── command/       # Domain commands
│   │   ├── entity/        # Domain entities
│   │   ├── event/         # Domain events
│   │   ├── repository/    # Repository interfaces
│   │   └── value/         # Value objects
│   ├── usecase/           # UseCase layer (application services)
│   │   ├── input/         # Input DTOs
│   │   └── output/        # Output DTOs
│   └── infrastructure/    # Infrastructure layer
│       ├── database/      # Database implementations
│       ├── handler/       # HTTP handlers
│       └── router/        # HTTP routing
├── container/             # Dependency injection container
├── config/               # Configuration management
├── docker-compose.yaml   # Docker services
└── Taskfile.yaml        # Task definitions
```

---

## License

This project is for educational purposes demonstrating Event Sourcing patterns in Go.