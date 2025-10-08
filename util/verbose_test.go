package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerbose(t *testing.T) {
	env_file.initialized = false

	assert.False(t,Verbose())

	verbose = nil
	env_file.initialized = true
	env_file.Content = map[string]string{
		"VERBOSE":"true",
	}

	assert.True(t,Verbose())
}