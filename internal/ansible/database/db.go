package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Ansible hosts database
type Database struct {
	dbFile string
	groups map[string]Group
}

// Create a new database
func NewDatabase(path string) *Database {
	return &Database{
		dbFile: fmt.Sprintf("%s%sterraform-provider-ansible.json", path, string(os.PathSeparator)),
		groups: make(map[string]Group),
	}
}

// Check if the database exists
func (s *Database) Exists() bool {
	_, err := os.Stat(s.dbFile)
	return err == nil
}

// Path to the database file
func (s *Database) Path() string {
	return s.dbFile
}

// Add a new ansible group to the database
func (s *Database) AddGroup(group Group) error {
	if _, ok := s.groups[group.GetID()]; ok {
		return fmt.Errorf("group '%s' already exists", group.GetID())
	}

	s.groups[group.GetID()] = group
	return nil
}

// Update an existing ansible group in the database
func (s *Database) UpdateGroup(group Group) {
	s.groups[group.GetID()] = group
}

// Remove an existing ansible group from the database
func (s *Database) RemoveGroup(group Group) error {
	if _, ok := s.groups[group.GetID()]; !ok {
		return nil
	}

	delete(s.groups, group.GetID())
	return nil
}

// Find a group with the specified ID in the database
func (s *Database) Group(id string) *Group {
	if val, ok := s.groups[id]; !ok {
		return nil
	} else {
		return &val
	}
}

// Find a host entry in the database by its ID and return the entry and which group it belongs to
func (s *Database) FindEntryById(id string) (*Group, *Entity, error) {
	for _, v := range s.groups {
		if e := v.Entry(id); e != nil {
			return &v, e, nil
		}
	}
	return nil, nil, fmt.Errorf("entry with GetID '%s' could not be found", id)
}

// Find a group in the database by its name
func (s *Database) FindGroupByName(name string) (*Group, error) {
	for k, v := range s.groups {
		if v.GetName() == name {
			g := s.groups[k]
			return &g, nil
		}
	}
	return nil, fmt.Errorf("group with name '%s' could not be found", name)
}

// Get a map of all the groups in the database
func (s *Database) AllGroups() *map[string]Group {
	return &s.groups
}

// Save the current in-memory version of the databse to disk
func (s *Database) Commit() error {
	// Commit JSON to disk
	if jsonString, err := json.MarshalIndent(s.groups, "", "\t"); err != nil {
		return fmt.Errorf("failed to serialize database to '%s': %e", s.dbFile, err)
	} else {
		if err := ioutil.WriteFile(s.dbFile, jsonString, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write database file '%s': %e", s.dbFile, err)
		}
	}

	return nil
}

// Load the database from disk into memory
func (s *Database) Load() error {
	if _, err := os.Stat(s.dbFile); os.IsNotExist(err) {
		return nil
	}

	s.groups = map[string]Group{}
	jsonString, err := ioutil.ReadFile(s.dbFile)
	if err != nil {
		return fmt.Errorf("failed to load database file '%s': %e", s.dbFile, err)
	}

	if len(jsonString) == 0 {
		return nil
	}

	if err := json.Unmarshal(jsonString, &s.groups); err != nil {
		return fmt.Errorf("failed to deserialize database '%s' to json: %e", s.dbFile, err)
	}

	return nil
}
