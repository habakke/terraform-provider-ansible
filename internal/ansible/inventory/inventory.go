package inventory

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Inventory represents an Ansible inventory
type Inventory struct {
	path          string
	fullPath      string
	groupVarsFile string
}

// Exists checks if an inventory with the given ID already exists
func Exists(id string) bool {
	_, err := os.Stat(id)
	return err == nil
}

// NewInventory creates a new inventory at the given path
func NewInventory(path string) Inventory {
	id := database.NewIdentity().GetID()
	return Inventory{
		path:          filepath.Clean(path),
		fullPath:      fmt.Sprintf("%s%s%s", filepath.Clean(path), string(os.PathSeparator), id),
		groupVarsFile: fmt.Sprintf("%s%s%s%sgroup_vars%sall.yml", path, string(os.PathSeparator), id, string(os.PathSeparator), string(os.PathSeparator)),
	}
}

// LoadFromID loads an inventory from disk given its ID
func LoadFromID(id string) Inventory {
	parts := strings.Split(filepath.Clean(id), string(os.PathSeparator))
	path := strings.Join(parts[:len(parts)-1], string(os.PathSeparator))
	return Inventory{
		path:          path,
		fullPath:      filepath.Clean(id),
		groupVarsFile: fmt.Sprintf("%s%sgroup_vars%sall.yml", filepath.Clean(id), string(os.PathSeparator), string(os.PathSeparator)),
	}
}

// GetID returns the ID of the inventory
func (s *Inventory) GetID() string {
	return s.fullPath
}

// GetPath returns the path to the inventory
func (s *Inventory) GetPath() string {
	return s.path
}

// GetDatabasePath returns the inventory database path
func (s *Inventory) GetDatabasePath() string {
	return s.fullPath
}

// GetAndLoadDatabase creates a new database and loads data from disk
func (s *Inventory) GetAndLoadDatabase() (*database.Database, error) {
	db := database.NewDatabase(s.GetDatabasePath())
	err := db.Load()
	return db, err
}

// Load loads the inventory from disk
func (s *Inventory) Load() (string, error) {
	if _, err := os.Stat(s.groupVarsFile); os.IsNotExist(err) {
		return "", fmt.Errorf("failed to load inventory because groupvars files doesn't exist: %s", err.Error())
	}

	if data, err := ioutil.ReadFile(s.groupVarsFile); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

// Commit saves groupVars for the inventory to disk
func (s *Inventory) Commit(groupVars string) error {
	if err := os.MkdirAll(fmt.Sprintf("%s%sgroup_vars", s.fullPath, string(os.PathSeparator)), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create inventory path: %s", err.Error())
	}

	if err := ioutil.WriteFile(s.groupVarsFile, []byte(groupVars), os.ModePerm); err != nil {
		return fmt.Errorf("failed to commit inventory to file: %s", err.Error())
	}

	return nil
}

// Delete deletes the inventory
func (s *Inventory) Delete() error {
	if _, err := os.Stat(s.fullPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(s.fullPath); err != nil {
		return fmt.Errorf("failed to delete inventory: %s", err.Error())
	}

	return nil
}
