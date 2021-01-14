package database

import "encoding/json"

type Host struct {
	id   Identity
	name string
}

func NewHost(name string) *Host {
	return &Host{
		id:   *NewIdentity(),
		name: name,
	}
}

func (s *Host) GetId() string {
	return s.id.GetId()
}

func (s *Host) GetName() string {
	return s.name
}

func (s *Host) SetName(name string) {
	s.name = name
}

func (s *Host) Type() string {
	return "HOST"
}

func (s Host) MarshalJSON() ([]byte, error) {
	aux := &struct {
		Id   Identity `json:"id"`
		Type string   `json:"type"`
		Name string   `json:"name"`
	}{
		Id:   s.id,
		Type: s.Type(),
		Name: s.name,
	}

	if jsonString, err := json.MarshalIndent(aux, "", "\t"); err != nil {
		return nil, err
	} else {
		return jsonString, err
	}
}

func (s *Host) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Id   Identity `json:"id"`
		Type string   `json:"type"`
		Name string   `json:"name"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.id = aux.Id
	s.name = aux.Name

	return nil
}
