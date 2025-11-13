package util

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type jt struct {
	Time   *JsonTime `json:"Time,omitempty"`
	TimeMs *JsonTimeMs `json:"TimeMs,omitempty"`
}

func TestJsonTime(t *testing.T) {
	dest := jt{}
	expected, err := time.Parse(JSONTIME_FORMAT, "2025-10-16T15:04:05Z")
	assert.Nil(t, err)
	err = json.Unmarshal([]byte(`{"Time":"2025-10-16T15:04:05Z"}`), &dest)
	assert.Nil(t, err)
	assert.Equal(t, JsonTime(expected), *dest.Time)
	j, _ := json.Marshal(dest)
	assert.Equal(t, `{"Time":"2025-10-16T15:04:05Z"}`, string(j))

	expected, err = time.Parse(JSONTIME_FORMAT, "2025-10-16T15:04:05.2345Z")
	assert.Nil(t, err)
	err = json.Unmarshal([]byte(`{"Time":"2025-10-16T15:04:05.2345Z"}`), &dest)
	assert.Nil(t, err)
	assert.Equal(t, JsonTime(expected), *dest.Time)
	j, _ = json.Marshal(dest)
	assert.Equal(t, `{"Time":"2025-10-16T15:04:05Z"}`, string(j))

	expected, err = time.Parse(JSONTIME_FORMAT, "2025-10-16T15:04:05.2345Z")
	assert.Nil(t, err)
	err = json.Unmarshal([]byte(`{"Time":"2025-10-16T15:04:05.2345Z"}`), &dest)
	assert.Nil(t, err)
	assert.Equal(t, JsonTime(expected), *dest.Time)
	j, _ = json.Marshal(dest)
	assert.Equal(t, `{"Time":"2025-10-16T15:04:05Z"}`, string(j))
}

func TestJsonTimeMs(t *testing.T) {
	dest := jt{}
	expected, err := time.Parse(JSONTIME_FORMAT, "2025-10-16T15:04:05Z")
	assert.Nil(t, err)
	err = json.Unmarshal([]byte(`{"TimeMs":"2025-10-16T15:04:05Z"}`), &dest)
	assert.Nil(t, err)
	assert.Equal(t, JsonTimeMs(expected), *dest.TimeMs)
	j, _ := json.Marshal(dest)
	assert.Equal(t, `{"TimeMs":"2025-10-16T15:04:05.000Z"}`, string(j))

	expected, err = time.Parse(JSONTIMEMS_FORMAT, "2025-10-16T15:04:05.234Z")
	assert.Nil(t, err)
	err = json.Unmarshal([]byte(`{"TimeMs":"2025-10-16T15:04:05.234Z"}`), &dest)
	assert.Nil(t, err)
	assert.Equal(t, JsonTimeMs(expected), *dest.TimeMs)
	j, _ = json.Marshal(dest)
	assert.Equal(t, `{"TimeMs":"2025-10-16T15:04:05.234Z"}`, string(j))

	expected, err = time.Parse(JSONTIMEMS_FORMAT, "2025-10-16T15:04:05.234Z")
	assert.Nil(t, err)
	err = json.Unmarshal([]byte(`{"TimeMs":"2025-10-16T15:04:05.234Z"}`), &dest)
	assert.Nil(t, err)
	assert.Equal(t, JsonTimeMs(expected), *dest.TimeMs)
	j, _ = json.Marshal(dest)
	assert.Equal(t, `{"TimeMs":"2025-10-16T15:04:05.234Z"}`, string(j))
}
