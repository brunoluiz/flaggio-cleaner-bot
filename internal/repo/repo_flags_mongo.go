package repo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FlagMongo defines mongo repository
type FlagMongo struct {
	db *mongo.Database
}

// NewFlagMongo returns a mongo repository
func NewFlagMongo(db *mongo.Database) *FlagMongo {
	return &FlagMongo{db}
}

// FindFlagsByMaxAge Find flags that are not updated for more than duration param
func (r *FlagMongo) FindFlagsByMaxAge(ctx context.Context, maxAge time.Duration) ([]Flag, error) {
	outdatedAt := time.Now().Add(maxAge * -1)
	cursor, err := r.db.Collection("flags").Find(ctx, bson.M{
		"updatedAt": bson.M{
			"$lt": outdatedAt,
		},
	})
	if err != nil {
		return []Flag{}, err
	}

	var flags []Flag
	if err = cursor.All(ctx, &flags); err != nil {
		return []Flag{}, err
	}

	return flags, nil
}
