package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentityValues(t *testing.T) {

	i1 := NewIdentity().GetID()
	i2 := NewIdentity().GetID()

	assert.NotEqual(t, i1, i2)
}
