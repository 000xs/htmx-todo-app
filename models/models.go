package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// User struct represents a user model
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

// Todo struct represents a todo model
type Todo struct {
	Id        uuid.UUID `json:"id"`
	Task      string    `json:"task"`
	Status    string    `json:"status"`
	UserID    string    `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}
