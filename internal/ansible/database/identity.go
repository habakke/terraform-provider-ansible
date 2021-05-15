package database

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
)

type Identity struct {
	id string
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
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
