package repo

import (
	"context"
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

const processedFlagsBoltBucket = "processed_flags"

// ProcessedFlagBolt defines a badger repository
type ProcessedFlagBolt struct {
	db *bolt.DB
}

// ProcessedFlagBolt returns a badger repository
func NewProcessedFlagDisk(path string) (*ProcessedFlagBolt, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	// Guarantee buckets
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(processedFlagsBoltBucket))
		if err != nil {
			return err
		}

		return nil
	})

	return &ProcessedFlagBolt{
		db: db,
	}, nil
}

// IsProcessed Checks if flag has been marked as outdated
func (r *ProcessedFlagBolt) IsProcessed(ctx context.Context, id string) (processed bool, err error) {
	err = r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(processedFlagsBoltBucket))
		if b == nil {
			processed = false
			return nil
		}

		if b.Get([]byte(id)) == nil {
			processed = false
			return nil
		}

		processed = true
		return nil
	})

	return processed, err
}

// Save Save details of resource
func (r *ProcessedFlagBolt) Save(ctx context.Context, flag ProcessedFlag) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(processedFlagsBoltBucket))

		buf, err := json.Marshal(flag)
		if err != nil {
			return err
		}

		return b.Put([]byte(flag.ID), buf)
	})
}

// Close Run garbage collection and close badger
func (r *ProcessedFlagBolt) Close() error {
	return r.db.Close()
}
