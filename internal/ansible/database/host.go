package database

import "encoding/json"

// Host represents an Ansible host in the hosts.ini file
type Host struct {
	id   Identity
	name string
}

// NewHost creates a new Host with the given name, where name is IP or hostname
func NewHost(name string) *Host {
	return &Host{
		id:   *NewIdentity(),
		name: name,
	}
}

// GetID returns the ID of the Host
func (s *Host) GetID() string {
	return s.id.GetID()
}

// GetName returns the name of the Host
func (s *Host) GetName() string {
	return s.name
}

// SetName sets the name of the Host
func (s *Host) SetName(name string) {
	s.name = name
}

// Type returns the Entity type of the Host
func (s *Host) Type() string {
	return "HOST"
}

// MarshalJSON marshals an Host to a JSON byte array
func (s Host) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ID   Identity `json:"id"`
		Type string   `json:"type"`
		Name string   `json:"name"`
	}{
		ID:   s.id,
		Type: s.Type(),
		Name: s.name,
	}

	if jsonString, err := json.MarshalIndent(aux, "", "\t"); err != nil {
		return nil, err
	} else {
		return jsonString, err
	}
}

// UnmarshalJSON returns an Host from a JSON byte array
func (s *Host) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID   Identity `json:"id"`
		Type string   `json:"type"`
		Name string   `json:"name"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.id = aux.ID
	s.name = aux.Name

	return nil
}
