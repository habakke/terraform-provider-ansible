package inventory

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const INVENTORY_PATH = "/tmp"

const TEST_GROUP_VARS_DATA = `---
k3s_version: v1.19.5+k3s1
ansible_user: ubuntu
systemd_dir: /etc/systemd/system
master_ip: "{{ hostvars[groups['master'][0]]['ansible_host'] | default(groups['master'][0]) }}"
extra_server_args: ""
extra_agent_args: ""
`

func TestCreateNewInventory(t *testing.T) {

	// 1. Create inventory

	i := NewInventory(INVENTORY_PATH)
	assert.Equal(t, INVENTORY_PATH, i.GetPath())
	if err := i.Commit(TEST_GROUP_VARS_DATA); err != nil {
		assert.Fail(t, "failed to create inventory")
	}

	if _, err := os.Stat(i.groupVarsFile); os.IsNotExist(err) {
		assert.Fail(t, "inventory group_vars file does not exist")
	}

	// 2. Load inventory from ID

	i2 := LoadFromId(i.GetId())
	assert.Equal(t, INVENTORY_PATH, i2.GetPath())
	if _, err := os.Stat(i2.groupVarsFile); os.IsNotExist(err) {
		assert.Fail(t, "inventory group_vars file does not exist")
	}

	if err, groupVars := i.Load(); err != nil {
		assert.Fail(t, "failed to load inventory")
	} else {
		assert.Equal(t, TEST_GROUP_VARS_DATA, groupVars)
	}

	// 3. Delete inventory

	if err := i.Delete(); err != nil {
		assert.Fail(t, "failed to delete inventory")
	}

	if _, err := os.Stat(i2.groupVarsFile); os.IsExist(err) {
		assert.Fail(t, "inventory exists even after delete")
	}

}
