package database

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Identity struct {
	id string
}

func generateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func NewIdentity() *Identity {
	rand.Seed(time.Now().UnixNano())
	return &Identity{
		id: generateId(),
	}
}

func (s *Identity) GetId() string {
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
