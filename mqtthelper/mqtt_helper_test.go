package mqtthelper

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestTopicMatch(t *testing.T) {
	assert.True(t,match("demo/get","demo/get"))
	assert.True(t,match("demo/get/aap","demo/get/aap"))
	assert.True(t,match("demo/get/aap","demo/get/#"))
	assert.True(t,match("demo/get/aap","demo/#"))
	assert.True(t,match("demo/get/question/aap","demo/get/+/aap"))
	assert.True(t,match("demo/get/question/aap","demo/get/+"))
}