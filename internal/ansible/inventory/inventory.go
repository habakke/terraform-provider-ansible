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
	id            string
	rootPath      string
	fullPath      string
	groupVarsFile string
}

// Exists checks if an inventory with the given ID already exists
func Exists(rootPath string, id string) bool {
	_, err := os.Stat(getFullPath(rootPath, id))
	return err == nil
}

// NewInventory creates a new inventory at the given rootPath
func NewInventory(rootPath string) Inventory {
	id := fmt.Sprintf("inventory%s%s", string(os.PathSeparator), database.NewIdentity().GetID())
	return Inventory{
		id:            id,
		rootPath:      filepath.Clean(rootPath),
		fullPath:      fmt.Sprintf("%s%s%s", filepath.Clean(rootPath), string(os.PathSeparator), id),
		groupVarsFile: fmt.Sprintf("%s%sall.yml", GetGroupVarsPath(rootPath, "all"), string(os.PathSeparator)),
	}
}

func checkId(id string) error {
	parts := strings.Split(filepath.Clean(id), string(os.PathSeparator))
	if parts[0] != "inventory" {
		return fmt.Errorf("invalid id, missing inventory part")
	}
	if len(parts) != 2 {
		return fmt.Errorf("incorrect id, does not contain enough rootPath elements")
	}
	return nil
}

func getFullPath(rootPath string, id string) string {
	return fmt.Sprintf("%s%s%s", filepath.Clean(rootPath), string(os.PathSeparator), id)
}

// Load loads an inventory from disk
func Load(rootPath string, id string) (*Inventory, error) {
	if err := checkId(id); err != nil {
		return nil, err
	}
	if _, err := os.Stat(getFullPath(rootPath, id)); os.IsExist(err) {
		return nil, fmt.Errorf("inventory not found at the expected location '%s': %s", getFullPath(rootPath, id), err.Error())
	}
	return &Inventory{
		id:            id,
		rootPath:      filepath.Clean(rootPath),
		fullPath:      getFullPath(rootPath, id),
		groupVarsFile: fmt.Sprintf("%s%sall.yml", GetGroupVarsPath(rootPath, "all"), string(os.PathSeparator)),
	}, nil
}

// GetID returns the ID of the inventory
func (s *Inventory) GetID() string {
	return s.id
}

// GetRootPath returns the root path where the inventories are stored
func (s *Inventory) GetRootPath() string {
	return s.rootPath
}

// GetInventoryPath returns the path to this specific inventory
func (s *Inventory) GetInventoryPath() string {
	return s.fullPath
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

	if data, err := ioutil.ReadFile(s.groupVarsFile); err != nil {
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

	if err := ioutil.WriteFile(s.groupVarsFile, []byte(groupVars), os.ModePerm); err != nil {
		return fmt.Errorf("failed to commit inventory to file: %s", err.Error())
	}

	return nil
}

func (s Inventory) getInventoryBasePath() string {
	return fmt.Sprintf("%s%sinventory", s.rootPath, string(os.PathSeparator))
}

func (s Inventory) deleteInventory() error {
	if _, err := os.Stat(s.getInventoryBasePath()); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(s.getInventoryBasePath()); err != nil {
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
