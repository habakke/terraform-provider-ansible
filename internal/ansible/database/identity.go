package database

import (
	"encoding/json"
	"fmt"
	"math/rand"
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
	return &Identity{
		id: generateId(),
	}
}

func NewIdentityFromId(id string) *Identity {
	return &Identity{
		id: id,
	}
}

func (s *Identity) GetId() string {
	if s.id == "" {
		s.id = generateId()
	}
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
