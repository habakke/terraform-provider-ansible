package ansible

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

const DB_PATH = "/tmp"
const ENCODE_FILE = "/tmp/encode_test.ini"

const TEST_HOST_DATA = `[master]
192.168.0.180
[node]
192.168.0.181
192.168.0.182
192.168.0.183
192.168.0.184
192.168.0.185
[k3s_cluster:children]
master
node
`

func TestExport(t *testing.T) {
	db := database.NewDatabase(DB_PATH)

	// add some test data
	master := database.NewGroup("master")
	_ = master.AddEntity(database.NewHost("192.168.0.180"))
	_ = db.AddGroup(*master)

	node := database.NewGroup("node")
	_ = node.AddEntity(database.NewHost("192.168.0.181"))
	_ = node.AddEntity(database.NewHost("192.168.0.182"))
	_ = node.AddEntity(database.NewHost("192.168.0.183"))
	_ = node.AddEntity(database.NewHost("192.168.0.184"))
	_ = node.AddEntity(database.NewHost("192.168.0.185"))
	_ = db.AddGroup(*node)

	groupInGroup := database.NewGroup("k3s_cluster:children")
	_ = groupInGroup.AddEntity(database.NewGroup("master"))
	_ = groupInGroup.AddEntity(database.NewGroup("node"))
	_ = db.AddGroup(*groupInGroup)

	// run test
	if err := Encode(ENCODE_FILE, db); err != nil {
		assert.Fail(t, fmt.Sprintf("failed to encode file: %e", err))
	}

	if data, err := ioutil.ReadFile(ENCODE_FILE); err != nil {
		assert.Fail(t, fmt.Sprintf("failed read encoded file: %e", err))
	} else {
		assert.Equal(t, len(TEST_HOST_DATA), len(string(data)))
	}
}
