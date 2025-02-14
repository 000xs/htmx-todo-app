package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/000xs/htmx-todo-app/db"
	"github.com/000xs/htmx-todo-app/models"
	"github.com/go-chi/chi"
	uuid "github.com/gofrs/uuid"
)

// TodoGetHandler retrieves a sample todo item
func TodoGetHandler(w http.ResponseWriter, r *http.Request) {
	reqLog(r)
	userID := chi.URLParam(r, "userId")
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Connect to the Redis database
	db, ctx := db.Connect()

	// Use SCAN to get all keys starting with "todo:"
	iter := db.Scan(*ctx, 0, "todo:*", 0).Iterator()
	var todos []models.Todo

	// Iterate over the keys
	for iter.Next(*ctx) {
		// For each todo key, retrieve its value
		todoData, err := db.Get(*ctx, iter.Val()).Result()
		if err != nil {
			http.Error(w, "Error fetching todo from Redis", http.StatusInternalServerError)
			return
		}

		// Unmarshal the todo data from JSON

		var todo models.Todo
		err = json.Unmarshal([]byte(todoData), &todo)
		if err != nil {
			// Log the error and data for debugging

			http.Error(w, "Error decoding todo JSON", http.StatusInternalServerError)
			return
		}
		if todo.UserID == userID {
			todos = append(todos, todo)

		}
		// Add the todo to the slice
		// todos = append(todos, todo)
	}

	if err := iter.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error iterating through Redis: %v", err), http.StatusInternalServerError)
		return
	}

	// Marshal the todos slice to JSON
	res, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, "Error encoding todos to JSON", http.StatusInternalServerError)
		return
	}

	// Write response with all todos
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

type GetTodo struct {
	Task   string `json:"task"`
	UserID string `json:"userId"`
}

func TodoPostHandler(w http.ResponseWriter, r *http.Request) {
	reqLog(r)
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Check if body is empty
	if len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	// Unmarshal JSON body into GetTodo struct
	var reqTodo GetTodo
	err = json.Unmarshal(body, &reqTodo)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Check if required fields are empty
	if reqTodo.Task == "" || reqTodo.UserID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Generate UUID for the new todo
	todoID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Error generating UUID", http.StatusInternalServerError)
		return
	}

	// Create todo struct
	todo := models.Todo{
		Id:        todoID, // Ensure the ID is converted properly (you may need a custom conversion depending on the UUID length)
		Task:      reqTodo.Task,
		Status:    "new",
		UserID:    reqTodo.UserID,
		CreatedAt: time.Now(),
	}

	// Marshal the Todo struct to JSON
	res, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, "Error encoding Todo to JSON", http.StatusInternalServerError)
		return
	}

	// add todo to database
	db, ctx := db.Connect()

	key := fmt.Sprintf("todo:%s", todo.Id)
	err = db.Set(*ctx, key, res, 0).Err()
	if err != nil {
		http.Error(w, "Error saving Todo to database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Todo created successfully: \n %s ", res)
	fmt.Printf("%vTodo created successfully: \n %v", "\033[35m", "\033[35m\n")

}

type UpdateTodo struct {
	TodoID string `json:"todoId"`
	Status string `json:"status"`
}

// Update Todo Handler
func TodoUpdateHandler(w http.ResponseWriter, r *http.Request) {
	reqLog(r)
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Check if body is empty
	if len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	// Unmarshal JSON body into UpdateTodo struct
	var reqUpdate UpdateTodo
	err = json.Unmarshal(body, &reqUpdate)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Check if required fields are empty
	if reqUpdate.TodoID == "" || reqUpdate.Status == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Validate status value
	if reqUpdate.Status != "in_progress" && reqUpdate.Status != "completed" {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	// Fetch the todo from the database
	db, ctx := db.Connect()
	key := fmt.Sprintf("todo:%s", reqUpdate.TodoID)
	todoData, err := db.Get(*ctx, key).Result()
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// Unmarshal the fetched todo data
	var todo models.Todo
	err = json.Unmarshal([]byte(todoData), &todo)
	if err != nil {
		http.Error(w, "Error decoding Todo from database", http.StatusInternalServerError)
		return
	}

	// Update the status of the todo
	todo.Status = reqUpdate.Status

	// Marshal the updated todo struct to JSON
	res, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, "Error encoding updated Todo to JSON", http.StatusInternalServerError)
		return
	}

	// Update the todo in the database
	err = db.Set(*ctx, key, res, 0).Err()
	if err != nil {
		http.Error(w, "Error updating Todo in database", http.StatusInternalServerError)
		return
	}

	// Respond with the updated todo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Todo updated successfully: \n %s", res)
	fmt.Printf("%vTodo updated successfully: \n %v", "\033[35m", "\033[35m\n")
}
