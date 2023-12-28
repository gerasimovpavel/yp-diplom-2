package domain

import "github.com/google/uuid"

type Player struct {
	UserID  uuid.UUID `json:"user_id" db:"user_id"`
	Place   int64     `json:"place" db:"place"`
	Name    string    `json:"name" db:"name"`
	Points  int64     `json:"points" db:"points"`
	Games   int64     `json:"games" db:"games"`
	Kills   int64     `json:"kills" db:"kills"`
	Bubbles int64     `json:"bubbles" db:"Bubbles"`
	IsAdmin bool      `json:"isAdmin" db:"IsAdmin"`
}
