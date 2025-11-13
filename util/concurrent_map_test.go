package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func val(v int) *int {
	return &v
}

func TestConcurrentMap(t *testing.T) {
	m := NewConcurrentMap[string,*int]()
	assert.Equal(t,0,m.Size())
	assert.True(t,m.Empty())
	assert.Nil(t, m.Put("5",val(5)))
	v, found := m.Get("5")
	assert.True(t,found)
	assert.Equal(t,5, *v)

	v, found = m.Get("not found")
	assert.False(t,found)
	assert.Nil(t, v)

	assert.Equal(t, 5, *m.Remove("5"))
	v, found = m.Get("5")
	assert.False(t,found)
	assert.Nil(t, v)

	v = m.ComputeIfAbsent("5", func() *int { 
		return val(8) 
	})
	assert.Equal(t, 8, *v)
	v = m.ComputeIfAbsent("5", func() *int { 
		return val(5) 
	})
	assert.Equal(t, 8, *v)

	assert.Equal(t,1,m.Size())
	assert.False(t,m.Empty())

}