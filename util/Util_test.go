package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestToHostPort(t *testing.T) {
	var hp HostPort
	var err error
	hp, err = ToHostPort("www.google.nl:443")
	assert.Equal(t,"www.google.nl",hp.Host)
	assert.Equal(t,443,hp.Port)
	assert.Nil(t,err)

	hp, err = ToHostPort("www.google.nl")
	assert.NotNil(t,err)
	assert.Empty(t,hp.Host)

	hp, err = ToHostPort("www.google.nl:https")
	assert.NotNil(t,err)
	assert.Empty(t,hp.Host)

	hp, err = ToHostPort("[2a00:1450:400e:80c::2004]:443")
	assert.Nil(t,err)
	assert.Equal(t,"[2a00:1450:400e:80c::2004]",hp.Host)
	assert.Equal(t,443,hp.Port)
}

func TestUriToHostPort(t *testing.T) {
	var hp HostPort
	var err error
	hp, err = UriToHostPort("https://www.google.nl:443")
	assert.Equal(t,"www.google.nl",hp.Host)
	assert.Equal(t,443,hp.Port)
	assert.Nil(t,err)

	hp, err = UriToHostPort("http://www.google.nl")
	assert.Equal(t,"www.google.nl",hp.Host)
	assert.Equal(t,80,hp.Port)
	assert.Nil(t,err)

	hp, err = UriToHostPort("http://www.google.nl/metwaterachter")
	assert.Equal(t,"www.google.nl",hp.Host)
	assert.Equal(t,80,hp.Port)
	assert.Nil(t,err)

	hp, err = UriToHostPort("huh://www.google.nl")
	assert.NotNil(t,err)
	assert.Empty(t,hp.Host)
}