package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Database struct {
	dbFile string
	groups map[string]Group
}

func NewDatabase(path string) *Database {
	return &Database{
		dbFile: fmt.Sprintf("%s%sterraform-provider-ansible.json", path, string(os.PathSeparator)),
		groups: make(map[string]Group),
	}
}

func (s *Database) Exists() bool {
	_, err := os.Stat(s.dbFile)
	return err == nil
}

func (s *Database) Path() string {
	return s.dbFile
}

func (s *Database) AddGroup(group Group) error {
	if _, ok := s.groups[group.GetId()]; ok {
		return errors.New(fmt.Sprintf("group '%s' already exists", group.GetId()))
	}

	s.groups[group.GetId()] = group
	return nil
}

func (s *Database) UpdateGroup(group Group) {
	s.groups[group.GetId()] = group
}

func (s *Database) RemoveGroup(group Group) error {
	if _, ok := s.groups[group.GetId()]; !ok {
		return nil
	}

	delete(s.groups, group.GetId())
	return nil
}

func (s *Database) Group(id string) *Group {
	if val, ok := s.groups[id]; !ok {
		return nil
	} else {
		return &val
	}
}

func (s *Database) FindEntryById(id string) (*Group, *Entity, error) {
	for _, v := range s.groups {
		if e := v.Entry(id); e != nil {
			return &v, e, nil
		}
	}
	return nil, nil, errors.New(fmt.Sprintf("entry with GetId '%s' could not be found", id))
}

func (s *Database) FindGroupByName(name string) (*Group, error) {
	for k, v := range s.groups {
		if v.GetName() == name {
			g := s.groups[k]
			return &g, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("group with name '%s' could not be found", name))
}

func (s *Database) AllGroups() *map[string]Group {
	return &s.groups
}

func (s *Database) Commit() error {
	// Commit JSON to disk
	if jsonString, err := json.MarshalIndent(s.groups, "", "\t"); err != nil {
		return errors.New(fmt.Sprintf("failed to serialize database to '%s': %e", s.dbFile, err))
	} else {
		if err := ioutil.WriteFile(s.dbFile, jsonString, os.ModePerm); err != nil {
			return errors.New(fmt.Sprintf("failed to write database file '%s': %e", s.dbFile, err))
		}
	}

	return nil
}

func (s *Database) Load() error {
	if _, err := os.Stat(s.dbFile); os.IsNotExist(err) {
		return nil
	}

	s.groups = map[string]Group{}
	if jsonString, err := ioutil.ReadFile(s.dbFile); err != nil {
		return errors.New(fmt.Sprintf("failed to load database file '%s': %e", s.dbFile, err))
	} else {
		if err := json.Unmarshal(jsonString, &s.groups); err != nil {
			return errors.New(fmt.Sprintf("failed to deserialize database '%s' to json: %e", s.dbFile, err))
		}
	}

	return nil
}
