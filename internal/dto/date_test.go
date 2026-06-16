package dto

import (
	"encoding/json"
	"testing"
)

func TestDateUnmarshal(t *testing.T) {
	cases := map[string]bool{
		`"2025-03-16"`:           true,  // формат date из OpenAPI
		`"2025-03-16T14:00:00Z"`: true,  // полный RFC3339
		`"16.03.2025"`:           false, // неверный формат
	}
	for input, ok := range cases {
		var d Date
		err := json.Unmarshal([]byte(input), &d)
		if ok && err != nil {
			t.Errorf("%s: unexpected error %v", input, err)
		}
		if !ok && err == nil {
			t.Errorf("%s: expected error, got none", input)
		}
		if ok && d.Year() != 2025 {
			t.Errorf("%s: wrong year %d", input, d.Year())
		}
	}
}
