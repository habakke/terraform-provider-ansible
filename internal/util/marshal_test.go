package util

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SuccessStruct struct {
	ID   string
	Name string
}

func (s SuccessStruct) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SuccessStruct) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return nil
}

type FailureStruct struct {
	ID   string
	Name string
}

func TestCanMarshal(t *testing.T) {
	assert.True(t, CanMarshal(SuccessStruct{ID: "1", Name: "Success test"}))
	assert.False(t, CanMarshal(FailureStruct{ID: "2", Name: "Failure test"}))
}
