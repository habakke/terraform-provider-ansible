package ansible

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

const DbPath = "/tmp"
const EncodeFile = "/tmp/encode_test.ini"

const TestHostData = `[master]
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
	db := database.NewDatabase(DbPath)

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
	_ = groupInGroup.AddEntity(database.NewGroup(master.GetName()))
	_ = groupInGroup.AddEntity(database.NewGroup(node.GetName()))
	_ = db.AddGroup(*groupInGroup)

	// run test
	if err := Encode(EncodeFile, db); err != nil {
		assert.Fail(t, fmt.Sprintf("failed to encode file: %e", err))
	}

	if data, err := ioutil.ReadFile(EncodeFile); err != nil {
		assert.Fail(t, fmt.Sprintf("failed read encoded file: %e", err))
	} else {
		fmt.Print(string(data))
		assert.Equal(t, len(TestHostData), len(string(data)))
	}
}
