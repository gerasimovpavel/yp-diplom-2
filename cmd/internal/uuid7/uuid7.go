package uuid7

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UUID7Bin primitive.Binary

func New() (uuid.UUID, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}
	return uid, nil
}
