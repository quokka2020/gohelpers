package util

import (
	"encoding/json"
	
	"time"
)
type JsonTime time.Time
type JsonTimeMs time.Time

const (
	JSONTIME_FORMAT = "2006-01-02T15:04:05Z"
	JSONTIMEMS_FORMAT = "2006-01-02T15:04:05.000Z"
)

func (c JsonTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(c).UTC().Format(JSONTIME_FORMAT))
}

func (c *JsonTime) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil
	}
	t, err := time.Parse(JSONTIME_FORMAT, j)
	if err == nil {
		*c = JsonTime(t)
	}
	return err
}

func (c JsonTimeMs) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(c).UTC().Format(JSONTIMEMS_FORMAT))
}

func (c *JsonTimeMs) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil
	}
	t, err := time.Parse(JSONTIMEMS_FORMAT, j)
	if err == nil {
		*c = JsonTimeMs(t)
	}
	t, err = time.Parse(JSONTIME_FORMAT, j)
	if err == nil {
		*c = JsonTimeMs(t)
	}
	return err
}