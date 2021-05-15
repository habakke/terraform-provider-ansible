package database

import (
	"encoding/json"
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/util"
)

// Ansible group
type Group struct {
	id      Identity
	name    string
	entries map[string]*Entity
}

// Creates a new Group with the given name
func NewGroup(name string) *Group {
	return &Group{
		id:      *NewIdentity(),
		name:    name,
		entries: make(map[string]*Entity),
	}
}

// Returns the ID of the Group
func (s *Group) GetID() string {
	return s.id.GetID()
}

// Returns the name of the Group
func (s *Group) GetName() string {
	return s.name
}

// Sets the name of the Group
func (s *Group) SetName(name string) {
	s.name = name
}

// Returns the Entity type
func (s *Group) Type() string {
	return "GROUP"
}

// Adds an Entity to the Group
func (s *Group) AddEntity(entity Entity) error {
	if _, ok := s.entries[entity.GetID()]; ok {
		return fmt.Errorf("entity '%s' already exists in group", entity.GetID())
	}

	s.entries[entity.GetID()] = &entity
	return nil
}

// Updates an Entity in the Group
func (s *Group) UpdateEntity(entity Entity) {
	s.entries[entity.GetID()] = &entity
}

// Removes an Entity from the Group
func (s *Group) RemoveEntity(entity Entity) error {
	if _, ok := s.entries[entity.GetID()]; !ok {
		return nil
	}

	delete(s.entries, entity.GetID())
	return nil
}

// Returns an Entity from the Group given its ID
func (s *Group) Entry(id string) *Entity {
	return s.entries[id]
}

// Returns a list of all Entity names in the Group as a string array
func (s *Group) GetEntriesAsString() []string {
	var stringEntries []string
	for k := range s.entries {
		stringEntries = append(stringEntries, (*s.entries[k]).GetName())
	}
	return stringEntries
}

func entriesMapToStringMap(entries map[string]*Entity) map[string]string {
	stringMap := make(map[string]string)
	for k, v := range entries {
		if !util.CanMarshal(*v) {
			continue
		}

		if s, err := json.Marshal(v); err == nil {
			stringMap[k] = string(s)
		}
	}
	return stringMap
}

// Marshal Group to JSON
func (s Group) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ID      Identity          `json:"id"`
		Type    string            `json:"type"`
		Name    string            `json:"name"`
		Entries map[string]string `json:"entries"`
	}{
		ID:      s.id,
		Type:    s.Type(),
		Name:    s.name,
		Entries: entriesMapToStringMap(s.entries),
	}

	if jsonString, err := json.MarshalIndent(aux, "", "\t"); err != nil {
		return nil, err
	} else {
		return jsonString, err
	}
}

// Unmarshal Group from JSON
func (s *Group) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID      Identity          `json:"id"`
		Type    string            `json:"type"`
		Name    string            `json:"name"`
		Entries map[string]string `json:"entries"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.id = aux.ID
	s.name = aux.Name
	s.entries = make(map[string]*Entity)

	for _, v := range aux.Entries {
		typeAux := &struct {
			Type string
		}{}

		if err := json.Unmarshal([]byte(v), &typeAux); err != nil {
			return err
		}

		switch typeAux.Type {
		case "HOST":
			h := &Host{}
			if err := json.Unmarshal([]byte(v), h); err != nil {
				return err
			}
			_ = s.AddEntity(h)
		case "GROUP":
			g := &Group{}
			if err := json.Unmarshal([]byte(v), g); err != nil {
				return err
			}
			_ = s.AddEntity(g)
		}
	}

	return nil
}
