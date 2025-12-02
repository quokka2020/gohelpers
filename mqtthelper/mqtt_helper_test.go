package mqtthelper

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestTopicMatch(t *testing.T) {
	assert.True(t,match("demo","demo/get","get"))
	assert.True(t,match("demo","demo/get/aap","get/aap"))
	assert.True(t,match("demo","demo/get/aap","get/#"))
	assert.True(t,match("demo","demo/get/aap","#"))
	assert.True(t,match("demo","demo/get/question/aap","get/+/aap"))
	assert.True(t,match("demo","demo/get/question/aap","get/+"))
}