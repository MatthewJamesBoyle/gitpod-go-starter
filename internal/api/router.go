package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"go-starter/internal/database"
	"go-starter/internal/kafka"
)

// TodoHandler handles todo-related HTTP requests
type TodoHandler struct {
	db       *sql.DB
	producer *kafka.Producer
}

// NewRouter sets up the HTTP routing
func NewRouter(db *sql.DB, producer *kafka.Producer) http.Handler {
	router := mux.NewRouter()
	handler := &TodoHandler{db: db, producer: producer}

	// API routes
	router.HandleFunc("/api/todos", handler.GetTodos).Methods("GET")
	router.HandleFunc("/api/todos", handler.CreateTodo).Methods("POST")
	router.HandleFunc("/api/todos/{id}", handler.UpdateTodo).Methods("PATCH")
	router.HandleFunc("/api/todos/{id}", handler.DeleteTodo).Methods("DELETE")

	// Serve static files for the frontend
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static")))

	return router
}

// GetTodos returns all todos
func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := database.GetTodos(h.db)
	if err != nil {
		log.Printf("Error getting todos: %v", err)
		http.Error(w, "Failed to get todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo creates a new todo
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todo, err := database.AddTodo(h.db, input.Title)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	// Send Kafka message about the new todo
	ctx := context.Background()
	if err := h.producer.SendMessage(ctx, "todos", "todo_created", todo); err != nil {
		log.Printf("Error sending Kafka message: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo updates a todo's status
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Completed bool `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := database.UpdateTodoStatus(h.db, id, input.Completed); err != nil {
		log.Printf("Error updating todo: %v", err)
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	// Send Kafka message about the updated todo
	ctx := context.Background()
	if err := h.producer.SendMessage(ctx, "todos", "todo_updated", map[string]interface{}{
		"id":        id,
		"completed": input.Completed,
	}); err != nil {
		log.Printf("Error sending Kafka message: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTodo deletes a todo
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := database.DeleteTodo(h.db, id); err != nil {
		log.Printf("Error deleting todo: %v", err)
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	// Send Kafka message about the deleted todo
	ctx := context.Background()
	if err := h.producer.SendMessage(ctx, "todos", "todo_deleted", map[string]int{
		"id": id,
	}); err != nil {
		log.Printf("Error sending Kafka message: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
