package repo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Flag Flag described by flaggio
type Flag struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Key         string             `json:"key"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ProcessedFlag Processed flag which should be cleaned-up
type ProcessedFlag struct {
	ID  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Key string             `json:"key"`
}
