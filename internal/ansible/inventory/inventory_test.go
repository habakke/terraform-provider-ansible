package inventory

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const InventoryPath = "/tmp"

const TestGroupVarsData = `---
k3s_version: v1.19.5+k3s1
ansible_user: ubuntu
systemd_dir: /etc/systemd/system
master_ip: "{{ hostvars[groups['master'][0]]['ansible_host'] | default(groups['master'][0]) }}"
extra_server_args: ""
extra_agent_args: ""
`

func TestCreateNewInventory(t *testing.T) {

	// 1. Create inventory

	i := NewInventory(InventoryPath)
	assert.Equal(t, InventoryPath, i.GetPath())
	if err := i.Commit(TestGroupVarsData); err != nil {
		assert.Fail(t, "failed to create inventory")
	}

	if _, err := os.Stat(i.groupVarsFile); os.IsNotExist(err) {
		assert.Fail(t, "inventory group_vars file does not exist")
	}

	// 2. Load inventory from ID

	i2 := LoadFromID(i.GetID())
	assert.Equal(t, InventoryPath, i2.GetPath())
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

	if _, err := os.Stat(i2.groupVarsFile); os.IsExist(err) {
		assert.Fail(t, "inventory exists even after delete")
	}

}
