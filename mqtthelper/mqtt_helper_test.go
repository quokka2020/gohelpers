package mqtthelper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTopicMatch(t *testing.T) {
	assert.True(t, match("demo/get", "demo/get"))
	assert.True(t, match("demo/get/aap", "demo/get/aap"))
	assert.True(t, match("demo/get/aap", "demo/get/#"))
	assert.True(t, match("demo/get/aap", "demo/#"))
	assert.True(t, match("demo/get/question/aap", "demo/get/+/aap"))
	assert.True(t, match("demo/get/question/aap", "demo/get/+"))
}

func TestValueToMessage(t *testing.T) {
	assert.Equal(t, []byte("aa"), ValueToMessage("aa"))
	assert.Equal(t, []byte("1"), ValueToMessage(true))
	assert.Equal(t, []byte("123.456000"), ValueToMessage(float64(123.456)))
	assert.Nil(t, ValueToMessage(nil))
}
