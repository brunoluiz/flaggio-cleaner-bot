package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProcessedFlagMongo defines a mongo repository
type ProcessedFlagMongo struct {
	db *mongo.Database
}

// NewProcessedFlagMongo returns a mongo repository
func NewProcessedFlagMongo(db *mongo.Database) *ProcessedFlagMongo {
	return &ProcessedFlagMongo{db}
}

// IsProcessed Checks if flag has been marked as outdated
func (r *ProcessedFlagMongo) IsProcessed(ctx context.Context, id primitive.ObjectID) (bool, error) {
	count, err := r.db.Collection("outdated").CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateProcessedFlag Mark flag as outdated
func (r *ProcessedFlagMongo) CreateProcessedFlag(ctx context.Context, flag ProcessedFlag) error {
	collection := r.db.Collection("outdated")

	_, err := collection.InsertOne(ctx, flag)
	if err != nil {
		return err
	}

	return nil
}
