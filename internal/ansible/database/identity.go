package database

import (
	"encoding/json"
	"github.com/google/uuid"
)

// Identity represents a database identity
type Identity struct {
	id string
}

func generateID() string {
	return uuid.New().String()
}

// NewIdentity creates a new database identity
func NewIdentity() *Identity {
	return &Identity{
		id: generateID(),
	}
}

// GetID returns the ID of an Identity
func (s *Identity) GetID() string {
	return s.id
}

// MarshalJSON marshals an Identity to a JSON byte array
func (s Identity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.id)
}

// UnmarshalJSON returns an Identity from a JSON byte array
func (s *Identity) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.id); err != nil {
		return err
	}
	return nil
}
