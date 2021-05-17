package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Database is an internal structure to represent the contents of an Ansible hosts.ini file
type Database struct {
	dbFile string
	groups map[string]Group
}

// NewDatabase creates a new database
func NewDatabase(path string) *Database {
	return &Database{
		dbFile: fmt.Sprintf("%s%sterraform-provider-ansible.json", path, string(os.PathSeparator)),
		groups: make(map[string]Group),
	}
}

// Exists checks if the database exists
func (s *Database) Exists() bool {
	_, err := os.Stat(s.dbFile)
	return err == nil
}

// Path to the database file
func (s *Database) Path() string {
	return s.dbFile
}

// AddGroup adds a new ansible group to the database
func (s *Database) AddGroup(group Group) error {
	if _, ok := s.groups[group.GetID()]; ok {
		return fmt.Errorf("group '%s' already exists", group.GetID())
	}

	s.groups[group.GetID()] = group
	return nil
}

// UpdateGroup updates an existing ansible group in the database
func (s *Database) UpdateGroup(group Group) {
	s.groups[group.GetID()] = group
}

// RemoveGroup removes an existing ansible group from the database
func (s *Database) RemoveGroup(group Group) error {
	if _, ok := s.groups[group.GetID()]; !ok {
		return nil
	}

	delete(s.groups, group.GetID())
	return nil
}

// Group returns a Group with the specified ID in the database
func (s *Database) Group(id string) *Group {
	if val, ok := s.groups[id]; !ok {
		return nil
	} else {
		return &val
	}
}

// FindEntryByID tries to locate a host entry in the database by its ID and return the entry and which Group it belongs to
func (s *Database) FindEntryByID(id string) (*Group, Entity, error) {
	for _, v := range s.groups {
		if e := v.Entry(id); e != nil {
			return &v, e, nil
		}
	}
	return nil, nil, fmt.Errorf("entry with GetID '%s' could not be found", id)
}

// FindGroupByName tries to locate a Group in the database by its name
func (s *Database) FindGroupByName(name string) (*Group, error) {
	for k, v := range s.groups {
		if v.GetName() == name {
			g := s.groups[k]
			return &g, nil
		}
	}
	return nil, fmt.Errorf("group with name '%s' could not be found", name)
}

// AllGroups returns a map of all the Groups in the database
func (s *Database) AllGroups() *map[string]Group {
	return &s.groups
}

// Commit the current in-memory version of the database to disk
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
