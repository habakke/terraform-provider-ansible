package database

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const DB_PATH = "/tmp"

func TestCreateNewDatabase(t *testing.T) {
	db := NewDatabase(DB_PATH)

	// add some test data
	master := NewGroup("master")
	_ = master.AddEntity(NewHost("192.168.0.180"))
	_ = db.AddGroup(*master)

	node := NewGroup("node")
	_ = node.AddEntity(NewHost("192.168.0.181"))
	_ = node.AddEntity(NewHost("192.168.0.182"))
	_ = node.AddEntity(NewHost("192.168.0.183"))
	_ = node.AddEntity(NewHost("192.168.0.184"))
	_ = node.AddEntity(NewHost("192.168.0.185"))
	_ = db.AddGroup(*node)

	groupInGroup := NewGroup("k3s_cluster:children")
	_ = groupInGroup.AddEntity(NewGroup("master"))
	_ = groupInGroup.AddEntity(NewGroup("node"))
	_ = db.AddGroup(*groupInGroup)

	if err := db.Commit(); err != nil {
		assert.Fail(t, fmt.Sprintf("%e", err))
	}

	// load the data back from disk again
	db2 := NewDatabase(DB_PATH)
	if err := db2.Load(); err != nil {
		assert.Fail(t, fmt.Sprintf("failed load db file: %e", err))
	}

	assert.Equal(t, 3, len(db2.groups))
	g1, err := db2.FindGroupByName("master")
	assert.Nil(t, err)
	assert.Equal(t, "master", g1.GetName())
	assert.Equal(t, 1, len(g1.entries))

	g2, err := db2.FindGroupByName("node")
	assert.Nil(t, err)
	assert.Equal(t, "node", g2.GetName())
	assert.Equal(t, 5, len(g2.entries))

	g3, err := db2.FindGroupByName("k3s_cluster:children")
	assert.Nil(t, err)
	assert.Equal(t, "k3s_cluster:children", g3.GetName())
	assert.Equal(t, 2, len(g3.entries))
}
