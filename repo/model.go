package repo

// ProcessedFlag Processed Flag details
type ProcessedFlag struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Key  string `json:"key"`
	Name string `json:"name"`
}
