package domain

import (
	"github.com/google/uuid"
	"time"
)

type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}

type User struct {
	UserID         uuid.UUID    `json:"user_id" db:"user_id" bson:"user_id"`
	Email          string       `json:"email" db:"email" bson:"email"`
	LastName       string       `json:"last_name" db:"last_name" bson:"last_name"`
	FirstName      string       `json:"first_name" db:"first_name" bson:"first_name"`
	Password       string       `json:"password" db:"password" bson:"password"`
	RegisteredAt   time.Time    `json:"registered_at" db:"registered_at" bson:"registered_at"`
	LastVisitAt    time.Time    `json:"last_visit_at" db:"last_visit_at" bson:"last_visit_at"`
	TelegramChatID int64        `json:"telegram_chat_id" db:"telegram_chat_id" bson:"telegram_chat_id"`
	Verification   Verification `json:"verification" db:"verification" bson:"verification"`
}

type UserDevices struct {
	UserID     uuid.UUID `json:"user_id" db:"user_id" bson:"user_id"`
	DeviceID   string    `json:"device_id" db:"device_id" bson:"device_id"`
	DeviceName string    `json:"device_name,omitempty" db:"device_name" bson:"device_name"`
	Current    bool      `json:"current" db:"current" bson:"current"`
}
