package database

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const DbPath = "/tmp"

func TestCreateNewDatabase(t *testing.T) {
	db := NewDatabase(DbPath)

	// add some test data
	master := NewGroup("master")
	variables := make(map[string]interface{})
	variables["test"] = "this is a test"
	_ = master.AddEntity(NewHost("192.168.0.180", variables))
	_ = db.AddGroup(*master)

	node := NewGroup("node")
	_ = node.AddEntity(NewHost("192.168.0.181", nil))
	_ = node.AddEntity(NewHost("192.168.0.182", nil))
	_ = node.AddEntity(NewHost("192.168.0.183", nil))
	_ = node.AddEntity(NewHost("192.168.0.184", nil))
	_ = node.AddEntity(NewHost("192.168.0.185", nil))
	_ = db.AddGroup(*node)

	groupInGroup := NewGroup("k3s_cluster:children")
	_ = groupInGroup.AddEntity(NewGroup("master"))
	_ = groupInGroup.AddEntity(NewGroup("node"))
	_ = db.AddGroup(*groupInGroup)

	if err := db.Commit(); err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err.Error()))
	}

	// load the data back from disk again
	db2 := NewDatabase(DbPath)
	if err := db2.Load(); err != nil {
		assert.Fail(t, fmt.Sprintf("failed load db file: %s", err.Error()))
	}

	assert.Equal(t, 3, len(db2.groups))
	g1, err := db2.FindGroupByName("master")
	assert.Nil(t, err)
	assert.Equal(t, "master", g1.GetName())
	assert.Equal(t, 1, len(g1.entries))

	e, err := g1.FindEntityByName("192.168.0.180")
	assert.Nil(t, err)
	h, ok := e.(*Host)
	assert.True(t, ok)
	assert.Equal(t, 1, len(h.variables))
	assert.Equal(t, "this is a test", h.variables["test"])

	g2, err := db2.FindGroupByName("node")
	assert.Nil(t, err)
	assert.Equal(t, "node", g2.GetName())
	assert.Equal(t, 5, len(g2.entries))

	g3, err := db2.FindGroupByName("k3s_cluster:children")
	assert.Nil(t, err)
	assert.Equal(t, "k3s_cluster:children", g3.GetName())
	assert.Equal(t, 2, len(g3.entries))
}
