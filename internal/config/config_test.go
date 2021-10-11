package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_valid(t *testing.T) {
	assert := require.New(t)
	t.Setenv("BACKEND_TABLE", "dynamodbtable")

	cfg := NewConfig()

	err := cfg.validate()
	assert.NoError(err)

	assert.Equal("dynamodbtable", cfg.BackendTable)
}
