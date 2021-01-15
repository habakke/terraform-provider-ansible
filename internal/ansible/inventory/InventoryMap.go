package inventory

import "fmt"

type InventoryMap struct {
	list map[string]Inventory
}

func NewInventoryMap() InventoryMap {
	return InventoryMap{
		list: make(map[string]Inventory),
	}
}

func (s *InventoryMap) Add(inventory Inventory) error {
	if _, ok := s.list[inventory.GetId()]; !ok {
		s.list[inventory.GetId()] = inventory
		return nil
	} else {
		return fmt.Errorf("inventory '%s' already exists", inventory.path)
	}
}

func (s *InventoryMap) Delete(id string) error {
	if _, ok := s.list[id]; ok {
		delete(s.list, id)
		return nil
	} else {
		return fmt.Errorf("inventory '%s' not found", id)
	}
}

func (s *InventoryMap) Get(id string) (*Inventory, error) {
	if i, ok := s.list[id]; ok {
		return &i, nil
	} else {
		return nil, fmt.Errorf("inventory '%s' not found", id)
	}
}

func (s *InventoryMap) Exists(id string) bool {
	_, ok := s.list[id]
	return ok
}

func (s *InventoryMap) GetOrCreate(id string) (*Inventory, error) {
	if i, ok := s.list[id]; ok {
		return &i, nil
	} else {
		i := LoadFromId(id)
		err := s.Add(i)
		return &i, err
	}
}
