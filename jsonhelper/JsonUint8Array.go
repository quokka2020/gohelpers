package jsonhelper

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JsonUint8Array []uint8

func (u JsonUint8Array) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}

func (u *JsonUint8Array) UnmarshalJSON(data []byte) error {
	var raw []uint8
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("JsonUint8Array: %w", err)
	}
	*u = raw
	return nil
}