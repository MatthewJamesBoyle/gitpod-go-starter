package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Connect establishes a connection to the PostgreSQL database
func Connect(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create the todos table if it doesn't exist
	if err := initializeSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// initializeSchema ensures the required tables exist
func initializeSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`

	_, err := db.Exec(query)
	return err
}

// GetTodos retrieves all todos from the database
func GetTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query("SELECT id, title, completed, created_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

// AddTodo inserts a new todo into the database
func AddTodo(db *sql.DB, title string) (Todo, error) {
	var todo Todo
	err := db.QueryRow(
		"INSERT INTO todos (title) VALUES ($1) RETURNING id, title, completed, created_at",
		title,
	).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt)

	return todo, err
}

// UpdateTodoStatus changes the completed status of a todo
func UpdateTodoStatus(db *sql.DB, id int, completed bool) error {
	result, err := db.Exec("UPDATE todos SET completed = $1 WHERE id = $2", completed, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with ID %d not found", id)
	}

	return nil
}

// DeleteTodo removes a todo from the database
func DeleteTodo(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with ID %d not found", id)
	}

	return nil
}

// Todo represents a task in the todo list
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}
