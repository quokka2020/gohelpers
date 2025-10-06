package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvFile(t *testing.T) {
	env_file.FileName = "env_file_test.env"
	env_file.initialized = false
	
	_,initialised := GetEnvFile()
	assert.False(t,initialised)
	user := os.Getenv("USER")
	assert.Equal(t,user,GetEnv("USER",""))

	_,initialised = GetEnvFile()
	assert.False(t,initialised)
	assert.Equal(t,5,GetEnvInt("FIVE",0))

	_,initialised = GetEnvFile()
	assert.True(t,initialised)

	assert.True(t,GetEnvBool("BOOL",false))
	assert.Equal(t,"asdf ! @",GetEnv("LINE",""))
	assert.Equal(t,"regel 1\nregel 2",GetEnv("MULTILINE",""))
	assert.Equal(t,"asdf ! @ # $ % ^ & * ( )",GetEnv("QUOTES",""))
	assert.Equal(t,fmt.Sprintf("%s asdf ! @ # $ %% ^ & * ( )",user),GetEnv("REPLACEMENT",""))
}