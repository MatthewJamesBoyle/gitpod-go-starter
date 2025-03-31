# Go Starter for Gitpod

A starter Go application with PostgreSQL and Kafka integration, ready for use in Gitpod.

## Features

- RESTful API for Todo management
- PostgreSQL database integration
- Kafka messaging for event-driven architecture
- Simple web frontend
- Development environment using Docker Compose
- Gitpod integration for cloud development

## Tech Stack

- Go (backend)
- PostgreSQL (database)
- Kafka (message broker)
- HTML/CSS/JavaScript (frontend)
- Docker Compose (containerization)
- Gitpod (cloud development)

## Getting Started

### Using Gitpod

The easiest way to get started is to use Gitpod:

1. Click the button below to open in Gitpod:

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/your-username/go-starter)

2. Gitpod will automatically start your development environment with all dependencies.
3. Your application will be available at port 8080.

### Local Development

To run the application locally:

1. Install Go 1.24 or later
2. Install Docker and Docker Compose
3. Clone the repository:
   ```
   git clone https://github.com/your-username/go-starter.git
   cd go-starter
   ```
4. Start the development environment:
   ```
   docker-compose up -d
   ```
5. Run the application:
   ```
   go run cmd/main.go
   ```
6. Open http://localhost:8080 in your browser

## Project Structure

```
├── cmd/                # Command line applications
│   └── main.go         # Main application entry point
├── internal/           # Private application code
│   ├── api/            # API handlers and routes
│   ├── config/         # Configuration management
│   ├── database/       # Database access layer
│   ├── kafka/          # Kafka integration
│   └── models/         # Domain models
├── web/                # Web assets
│   └── static/         # Static files for the frontend
├── .devcontainer.json  # VS Code development container configuration
├── docker-compose.yaml # Docker compose configuration
├── go.mod              # Go module definition
└── README.md           # This file
```

## API Endpoints

- `GET /api/todos` - List all todos
- `POST /api/todos` - Create a new todo
- `PATCH /api/todos/:id` - Update a todo's completion status
- `DELETE /api/todos/:id` - Delete a todo

## Environment Variables

The application can be configured using environment variables:

- `SERVER_PORT` - Port for the HTTP server (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string (default: postgres://todos:todos@db:5432/todos?sslmode=disable)
- `KAFKA_BROKERS` - Kafka broker address (default: kafka:9092)
- `KAFKA_TOPIC` - Kafka topic name (default: todos)

## License

MIT 