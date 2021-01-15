package inventory

import (
	"errors"
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Inventory struct {
	path          string
	fullPath      string
	groupVarsFile string
}

func Exists(id string) bool {
	_, err := os.Stat(id)
	return err == nil
}

func NewInventory(path string) Inventory {
	id := database.NewIdentity().GetId()
	return Inventory{
		path:          filepath.Clean(path),
		fullPath:      fmt.Sprintf("%s%s%s", filepath.Clean(path), string(os.PathSeparator), id),
		groupVarsFile: fmt.Sprintf("%s%s%s%sgroup_vars%sall.yml", path, string(os.PathSeparator), id, string(os.PathSeparator), string(os.PathSeparator)),
	}
}

func LoadFromId(id string) Inventory {
	parts := strings.Split(filepath.Clean(id), string(os.PathSeparator))
	path := strings.Join(parts[:len(parts)-1], string(os.PathSeparator))
	return Inventory{
		path:          path,
		fullPath:      filepath.Clean(id),
		groupVarsFile: fmt.Sprintf("%s%sgroup_vars%sall.yml", filepath.Clean(id), string(os.PathSeparator), string(os.PathSeparator)),
	}
}

func (s *Inventory) GetId() string {
	return s.fullPath
}

func (s *Inventory) GetPath() string {
	return s.path
}

func (s *Inventory) GetDatabasePath() string {
	return s.fullPath
}

func (s *Inventory) GetAndLoadDatabase() (*database.Database, error) {
	db := database.NewDatabase(s.GetDatabasePath())
	err := db.Load()
	return db, err
}

func (s *Inventory) Load() (error, string) {
	if _, err := os.Stat(s.groupVarsFile); os.IsNotExist(err) {
		return errors.New("failed to load inventory because groupvars files doesn't exist"), ""
	}

	if data, err := ioutil.ReadFile(s.groupVarsFile); err != nil {
		return err, ""
	} else {
		return nil, string(data)
	}
}

func (s *Inventory) Commit(groupVars string) error {
	if err := os.MkdirAll(fmt.Sprintf("%s%sgroup_vars", s.fullPath, string(os.PathSeparator)), os.ModePerm); err != nil {
		return errors.New("failed to create inventory path")
	}

	if err := ioutil.WriteFile(s.groupVarsFile, []byte(groupVars), os.ModePerm); err != nil {
		return errors.New("failed to commit inventory to file")
	}

	return nil
}

func (s *Inventory) Delete() error {
	if _, err := os.Stat(s.fullPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(s.fullPath); err != nil {
		return fmt.Errorf("failed to delete inventory: %e", err)
	}

	return nil
}
