package util

import (
	"encoding/json"
	
	"time"
)
type JsonTime time.Time

const JSONTIME_FORMAT = "2006-01-02T15:04:05Z"

func (c JsonTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(c).UTC().Format(JSONTIME_FORMAT))
}

func (c *JsonTime) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil
	}
	time, err := time.Parse(JSONTIME_FORMAT, j)
	if err == nil {
		*c = JsonTime(time)
	}
	return err
}