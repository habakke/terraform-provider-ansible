package inventory

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"os"
	"path/filepath"
)

// Inventory represents an Ansible inventory
type Inventory struct {
	id            string
	rootPath      string
	groupVarsFile string
}

// Exists checks if an inventory with the given ID already exists
func Exists(rootPath string, id string) bool {
	_, err := os.Stat(rootPath)
	if err != nil {
		return false
	}
	actualID, err := getId(rootPath)
	if err != nil {
		return false
	}
	if actualID != id {
		return false
	}
	return true
}

// NewInventory creates a new inventory at the given rootPath
func NewInventory(rootPath string) Inventory {
	return Inventory{
		id:            database.NewIdentity().GetID(),
		rootPath:      filepath.Clean(rootPath),
		groupVarsFile: fmt.Sprintf("%s%sall.yml", GetGroupVarsPath(rootPath, "all"), string(os.PathSeparator)),
	}
}

func writeId(rootPath string, id string) error {
	return os.WriteFile(fmt.Sprintf("%s/id", filepath.Clean(rootPath)), []byte(id), os.ModePerm)
}

func getId(rootPath string) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/id", filepath.Clean(rootPath)))
	if err != nil {
		return "", err
	}
	return string(data), err
}

// Load loads an inventory from disk
func Load(rootPath string, id string) (*Inventory, error) {
	actualID, err := getId(rootPath)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(rootPath); os.IsExist(err) {
		return nil, fmt.Errorf("inventory not found at the expected location '%s': %s", rootPath, err.Error())
	}
	if actualID != id {
		return nil, fmt.Errorf("inventory found, but ID is incorrect (id=%s, actualID=%s)", id, actualID)
	}
	return &Inventory{
		id:            id,
		rootPath:      filepath.Clean(rootPath),
		groupVarsFile: fmt.Sprintf("%s%sall.yml", GetGroupVarsPath(rootPath, "all"), string(os.PathSeparator)),
	}, nil
}

// GetID returns the ID of the inventory
func (s *Inventory) GetID() string {
	return s.id
}

// GetInventoryPath returns the path to this specific inventory
func (s *Inventory) GetInventoryPath() string {
	return s.rootPath
}

// GetGroupVarsPath returns the rootPath to the group_vars folder for an ansible group in the inventory
func GetGroupVarsPath(path string, group string) string {
	if len(group) == 0 {
		return fmt.Sprintf("%s%sgroup_vars", path, string(os.PathSeparator))

	} else {
		return fmt.Sprintf("%s%sgroup_vars%s%s", path, string(os.PathSeparator), string(os.PathSeparator), group)

	}
}

// GetAndLoadDatabase creates a new database and loads data from disk
func (s *Inventory) GetAndLoadDatabase() (*database.Database, error) {
	db := database.NewDatabase(s.GetInventoryPath())
	err := db.Load()
	return db, err
}

// Load loads the inventory from disk
func (s *Inventory) Load() (string, error) {
	if _, err := os.Stat(s.groupVarsFile); os.IsNotExist(err) {
		return "", fmt.Errorf("failed to load inventory because groupvars files doesn't exist: %s", err.Error())
	}

	if data, err := os.ReadFile(s.groupVarsFile); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

// Commit saves groupVars for the inventory to disk
func (s *Inventory) Commit(groupVars string) error {
	if err := os.MkdirAll(s.GetInventoryPath(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create inventory rootPath: %s", err.Error())
	}
	if err := os.MkdirAll(GetGroupVarsPath(s.rootPath, "all"), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create inventory group_vars rootPath: %s", err.Error())
	}
	if err := writeId(s.rootPath, s.id); err != nil {
		return fmt.Errorf("failed to write inventory id: %s", err.Error())
	}
	if err := os.WriteFile(s.groupVarsFile, []byte(groupVars), os.ModePerm); err != nil {
		return fmt.Errorf("failed to commit inventory to file: %s", err.Error())
	}
	return nil
}

func (s Inventory) deleteInventory() error {
	if _, err := os.Stat(s.rootPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(s.rootPath); err != nil {
		return fmt.Errorf("failed to delete inventory: %s", err.Error())
	}
	return nil
}

func (s Inventory) deleteGroupVars() error {
	if _, err := os.Stat(GetGroupVarsPath(s.rootPath, "")); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(GetGroupVarsPath(s.rootPath, "")); err != nil {
		return fmt.Errorf("failed to delete group_vars: %s", err.Error())
	}
	return nil
}

// Delete deletes the inventory
func (s *Inventory) Delete() error {
	if err := s.deleteGroupVars(); err != nil {
		return err
	}
	if err := s.deleteInventory(); err != nil {
		return err
	}
	return nil
}
