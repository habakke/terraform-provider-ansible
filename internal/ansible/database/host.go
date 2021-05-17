package database

import (
	"encoding/json"
	"fmt"
)

// Host represents an Ansible host in the hosts.ini file
type Host struct {
	id        Identity
	name      string
	variables map[string]interface{}
}

// NewHost creates a new Host with the given name, where name is IP or hostname
func NewHost(name string, variables map[string]interface{}) *Host {
	vars := variables
	if vars == nil {
		vars = make(map[string]interface{})
	}

	return &Host{
		id:        *NewIdentity(),
		name:      name,
		variables: vars,
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

// GetVariableNames returns the name of all variables set for a host
func (s *Host) GetVariableNames() []string {
	keys := make([]string, 0, len(s.variables))
	for k := range s.variables {
		keys = append(keys, k)
	}
	return keys
}

// GetVariables returns variable map
func (s *Host) GetVariables() map[string]interface{} {
	return s.variables
}

// GetVariable returns a variable for a host
func (s *Host) GetVariable(name string) (interface{}, error) {
	if val, ok := s.variables[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("variable '%s' not defined for host '%s'", name, s.name)
}

// SetVariable sets a variable for a host
func (s *Host) SetVariable(name string, val interface{}) {
	s.variables[name] = val
}

// Type returns the Entity type of the Host
func (s *Host) Type() string {
	return "HOST"
}

// MarshalJSON marshals an Host to a JSON byte array
func (s Host) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ID        Identity               `json:"id"`
		Type      string                 `json:"type"`
		Name      string                 `json:"name"`
		Variables map[string]interface{} `json:"variables"`
	}{
		ID:        s.id,
		Type:      s.Type(),
		Name:      s.name,
		Variables: s.variables,
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
		ID        Identity               `json:"id"`
		Type      string                 `json:"type"`
		Name      string                 `json:"name"`
		Variables map[string]interface{} `json:"variables"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.id = aux.ID
	s.name = aux.Name
	s.variables = aux.Variables

	return nil
}
