package flaggio

import "time"

// Flag Flag described by flaggio
type Flag struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

// FindFlagsResponse GQL findFlags response
type FindFlagsResponse struct {
	Flags struct {
		Flags []Flag `json:"flags"`
		Total int    `json:"total"`
	} `json:"flags"`
}
