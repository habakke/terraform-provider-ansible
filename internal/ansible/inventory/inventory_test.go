package inventory

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const InventoryRootPath = "/tmp"

const TestGroupVarsData = `---
k3s_version: v1.19.5+k3s1
ansible_user: ubuntu
systemd_dir: /etc/systemd/system
master_ip: "{{ hostvars[groups['master'][0]]['ansible_host'] | default(groups['master'][0]) }}"
extra_server_args: ""
extra_agent_args: ""
`

func TestInventoryPaths(t *testing.T) {
	i := NewInventory(InventoryRootPath)
	assert.Equal(t, InventoryRootPath, i.GetRootPath())
	assert.Equal(t, fmt.Sprintf("%s/inventory", InventoryRootPath), i.getInventoryBasePath())
	assert.Equal(t, fmt.Sprintf("%s%s%s", InventoryRootPath, string(os.PathSeparator), i.GetID()), i.GetInventoryPath())
	assert.Equal(t, i.id, i.GetID())
	assert.Equal(t, fmt.Sprintf("%s/group_vars/all", InventoryRootPath), GetGroupVarsPath(InventoryRootPath, "all"))
}

func TestCreateNewInventory(t *testing.T) {

	// 1. Create inventory

	i := NewInventory(InventoryRootPath)
	assert.Equal(t, InventoryRootPath, i.GetRootPath())
	if err := i.Commit(TestGroupVarsData); err != nil {
		assert.Fail(t, "failed to create inventory")
	}

	if _, err := os.Stat(i.groupVarsFile); os.IsNotExist(err) {
		assert.Fail(t, "inventory group_vars file does not exist")
	}

	// 2. Load inventory from ID

	i2, err := Load(i.GetRootPath(), i.GetID())
	assert.NoError(t, err)
	assert.Equal(t, InventoryRootPath, i2.GetRootPath())
	if _, err := os.Stat(i2.groupVarsFile); os.IsNotExist(err) {
		assert.Fail(t, "inventory group_vars file does not exist")
	}

	if groupVars, err := i.Load(); err != nil {
		assert.Fail(t, "failed to load inventory")
	} else {
		assert.Equal(t, TestGroupVarsData, groupVars)
	}

	// 3. Delete inventory

	if err := i.Delete(); err != nil {
		assert.Fail(t, "failed to delete inventory")
	}

	if _, err := os.Stat(i2.GetInventoryPath()); os.IsExist(err) {
		assert.Fail(t, "inventory exists even after delete")
	}
	if _, err := os.Stat(i2.groupVarsFile); os.IsExist(err) {
		assert.Fail(t, "inventory group_vars exists even after delete")
	}
}
