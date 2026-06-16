package dto

import (
	"fmt"
	"strings"
	"time"
)

// Date принимает дату как "2006-01-02" (формат `date` из OpenAPI) либо как
// полный RFC3339. Нужно, чтобы фронт мог слать event_date в виде "2025-03-16".
type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		d.Time = t
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("invalid date %q: expected YYYY-MM-DD", s)
	}
	d.Time = t
	return nil
}
