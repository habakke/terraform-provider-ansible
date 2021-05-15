package database

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Identity struct {
	id string
}

func generateID() string {
	return uuid.New().String()
}

func NewIdentity() *Identity {
	return &Identity{
		id: generateID(),
	}
}

func (s *Identity) GetID() string {
	return s.id
}

func (s Identity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.id)
}

func (s *Identity) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.id); err != nil {
		return err
	}
	return nil
}
