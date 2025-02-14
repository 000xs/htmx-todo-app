package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/000xs/htmx-todo-app/db"
	"github.com/000xs/htmx-todo-app/models"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

var users = map[string]string{}          // In-memory map for user data (username -> password)
var secretKey = []byte("supersecretkey") // Secret key for JWT signing

// User represents a user
type ReqUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Allusers(db *redis.Client, ctx *context.Context) []models.User {

	iter := db.Scan(*ctx, 0, "user:*", 0).Iterator()
	var users []models.User

	// Iterate over the keys
	for iter.Next(*ctx) {
		// For each todo key, retrieve its value
		userData, err := db.Get(*ctx, iter.Val()).Result()
		if err != nil {

			return nil
		}

		var user models.User
		err = json.Unmarshal([]byte(userData), &user)
		if err != nil {
			// Log the error and data for debugging

			return nil
		}

		users = append(users, user)

		// Add the todo to the slice
		// todos = append(todos, todo)
	}

	if err := iter.Err(); err != nil {

		return nil
	}

	return users
}

// RegisterHandler - Registers a new user
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	reqLog(r)
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
	var reqUser ReqUser
	err = json.Unmarshal(body, &reqUser)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Check if required fields are empty
	if reqUser.Username == "" || reqUser.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Generate UUID for the new todo
	userID, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Error generating UUID", http.StatusInternalServerError)
		fmt.Printf("Error generating UUID \n")
		return
	}

	user := models.User{
		userID,
		reqUser.Username,
		reqUser.Password,
	}

	// user struct convets to json
	jsonUser, err := json.Marshal(&user)

	if err != nil {
		http.Error(w, "Error encoding user to JSON", http.StatusInternalServerError)
		return
	}

	//add user to redis
	db, ctx := db.Connect()

	//fetch all users
	allusers := Allusers(db, ctx)
	if allusers == nil {
		http.Error(w, "Error fetching all users", http.StatusInternalServerError)
		return

	}

	for i := 0; i < len(allusers); i++ {
		if allusers[i].Username == user.Username {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}
	}

	key := fmt.Sprintf("user:%s", user.ID)
	err = db.Set(*ctx, key, jsonUser, 0).Err()
	if err != nil {
		http.Error(w, "Error saving Todo to database", http.StatusInternalServerError)
		return
	}
	// Generate JWT token for the user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.ID.String(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the JWT token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the response with the user details and JWT token
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "User created successfully",
		"user_id": user.ID.String(),
		"token":   tokenString,
	}
	json.NewEncoder(w).Encode(response)

}

// LoginHandler - Authenticates a user and returns a JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	reqLog(r)
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

	// Unmarshal JSON body into ReqUser struct
	var reqUser ReqUser
	err = json.Unmarshal(body, &reqUser)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Check if required fields are empty
	if reqUser.Username == "" || reqUser.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Fetch user data from Redis (for simplicity, only check username and password)
	db, ctx := db.Connect()
	allusers := Allusers(db, ctx)
	if allusers == nil {
		http.Error(w, "Error fetching all users", http.StatusInternalServerError)
		fmt.Printf("Error fetching all users \n")
		return
	}

	var user models.User
	for _, u := range allusers {
		if u.Username == reqUser.Username && u.Password == reqUser.Password {
			user = u
			break
		}
	}

	// If user not found or password doesn't match
	if user.ID == uuid.Nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token for the user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.ID.String(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the JWT token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the JWT token
	response := map[string]interface{}{
		"message": "Login successful",
		"user_id": user.ID.String(),
		"token":   tokenString,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ValidateHandler - Checks if the request has a valid JWT token
func ValidateHandler(w http.ResponseWriter, r *http.Request) {

	reqLog(r)
	// Get the token from Authorization header (e.g., "Bearer <token>")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing authorization header", http.StatusUnauthorized)
		return
	}

	// Parse the token (expecting format: "Bearer <token>")
	tokenString := authHeader[len("Bearer "):]

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Extract claims if token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		userID := claims["user_id"].(string)

		// Respond with the user information
		response := map[string]interface{}{
			"message":     "Token is valid",
			"username":    username,
			"user_id":     userID,
			"valid_until": claims["exp"],
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}
}
