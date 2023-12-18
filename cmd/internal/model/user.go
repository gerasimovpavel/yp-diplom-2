package model

import "github.com/google/uuid"

type User struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Login  string    `json:"login" db:"login"`
}
