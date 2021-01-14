package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/util"
)

type Group struct {
	id      Identity
	name    string
	entries map[string]*Entity
}

func NewGroup(name string) *Group {
	return &Group{
		id:      *NewIdentity(),
		name:    name,
		entries: make(map[string]*Entity),
	}
}

func (s *Group) GetId() string {
	return s.id.GetId()
}

func (s *Group) GetName() string {
	return s.name
}

func (s *Group) SetName(name string) {
	s.name = name
}

func (s *Group) Type() string {
	return "GROUP"
}

func (s *Group) AddEntity(entity Entity) error {
	if _, ok := s.entries[entity.GetId()]; ok {
		return errors.New(fmt.Sprintf("entity '%s' already exists in group", entity.GetId()))
	}

	s.entries[entity.GetId()] = &entity
	return nil
}

func (s *Group) UpdateEntity(entity Entity) {
	s.entries[entity.GetId()] = &entity
}

func (s *Group) RemoveEntity(entity Entity) error {
	if _, ok := s.entries[entity.GetId()]; !ok {
		return nil
	}

	delete(s.entries, entity.GetId())
	return nil
}

func (s *Group) Entry(id string) *Entity {
	return s.entries[id]
}

func (s *Group) GetEntriesAsString() []string {
	var stringEntries []string
	for k, _ := range s.entries {
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

func (s Group) MarshalJSON() ([]byte, error) {
	aux := &struct {
		Id      Identity          `json:"id"`
		Type    string            `json:"type"`
		Name    string            `json:"name"`
		Entries map[string]string `json:"entries"`
	}{
		Id:      s.id,
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

func (s *Group) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Id      Identity          `json:"id"`
		Type    string            `json:"type"`
		Name    string            `json:"name"`
		Entries map[string]string `json:"entries"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.id = aux.Id
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
			s.AddEntity(h)
		case "GROUP":
			g := &Group{}
			if err := json.Unmarshal([]byte(v), g); err != nil {
				return err
			}
			s.AddEntity(g)
		}
	}

	return nil
}
