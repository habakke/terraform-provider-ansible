package database

import "encoding/json"

// Ansible Host
type Host struct {
	id   Identity
	name string
}

// Create a new Host with the given name, where name is IP or hostname
func NewHost(name string) *Host {
	return &Host{
		id:   *NewIdentity(),
		name: name,
	}
}

// Returns the ID of the Host
func (s *Host) GetID() string {
	return s.id.GetId()
}

// Returns the name of the Host
func (s *Host) GetName() string {
	return s.name
}

// Sets the name of the Host
func (s *Host) SetName(name string) {
	s.name = name
}

// Returns the Entity type of the Host
func (s *Host) Type() string {
	return "HOST"
}

// Marshals the Host to JSON
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

// Unmarshal a Host from JSON
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
