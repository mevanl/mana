package util

import "encoding/json"

// MustMarshal marshals data and panics on failure.
// Use only when you're confident data is serializable.
func MustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic("mustMarshal failed: " + err.Error())
	}
	return b
}
