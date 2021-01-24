package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentityValues(t *testing.T) {

	i1 := NewIdentity().GetId()
	i2 := NewIdentity().GetId()

	assert.NotEqual(t, i1, i2)
}
