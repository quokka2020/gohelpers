package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerbose(t *testing.T) {
	verbose = -1
	env_file.initialized = true
	env_file.Content = map[string]string{
		"VERBOSE":"true",
	}

	assert.True(t,Verbose())
}

func TestVerboseLevel(t *testing.T) {
	verbose = -1

	env_file.initialized = true
	env_file.Content = map[string]string{
		"VERBOSE":"4",
	}

	assert.Equal(t,4,VerboseLevel())
	assert.True(t,Verbose())
}