package jsonhelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type JsonDuration time.Duration

func durationToISO8601(d time.Duration) string {
	if d == 0 {
		return "PT0S"
	}

	neg := d < 0
	if neg {
		d = -d
	}

	totalSeconds := int(d.Seconds())
	days := totalSeconds / 86400
	remainder := totalSeconds % 86400
	hours := remainder / 3600
	remainder %= 3600
	minutes := remainder / 60
	seconds := remainder % 60

	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	b.WriteByte('P')
	if days > 0 {
		fmt.Fprintf(&b, "%dD", days)
	}
	if hours > 0 || minutes > 0 || seconds > 0 {
		b.WriteByte('T')
		if hours > 0 {
			fmt.Fprintf(&b, "%dH", hours)
		}
		if minutes > 0 {
			fmt.Fprintf(&b, "%dM", minutes)
		}
		if seconds > 0 {
			fmt.Fprintf(&b, "%dS", seconds)
		}
	}
	return b.String()
}

func parseISO8601(s string) (time.Duration, error) {
	if len(s) == 0 {
		return 0, errors.New("empty ISO8601 duration")
	}

	neg := false
	rest := s
	if rest[0] == '-' {
		neg = true
		rest = rest[1:]
	}

	if len(rest) == 0 || rest[0] != 'P' {
		return 0, fmt.Errorf("invalid ISO8601 duration: %q", s)
	}
	rest = rest[1:]

	if len(rest) == 0 {
		return 0, fmt.Errorf("invalid ISO8601 duration: %q", s)
	}

	var d time.Duration
	inTime := false

	for len(rest) > 0 {
		if rest[0] == 'T' {
			inTime = true
			rest = rest[1:]
			if len(rest) == 0 {
				return 0, fmt.Errorf("invalid ISO8601 duration: %q", s)
			}
			continue
		}

		// Parse number
		n := 0
		found := false
		for len(rest) > 0 && rest[0] >= '0' && rest[0] <= '9' {
			n = n*10 + int(rest[0]-'0')
			rest = rest[1:]
			found = true
		}
		if !found || len(rest) == 0 {
			return 0, fmt.Errorf("invalid ISO8601 duration: %q", s)
		}

		unit := rest[0]
		rest = rest[1:]

		switch {
		case unit == 'D' && !inTime:
			d += time.Duration(n) * 24 * time.Hour
		case unit == 'H' && inTime:
			d += time.Duration(n) * time.Hour
		case unit == 'M' && inTime:
			d += time.Duration(n) * time.Minute
		case unit == 'S' && inTime:
			d += time.Duration(n) * time.Second
		default:
			return 0, fmt.Errorf("invalid ISO8601 duration: %q", s)
		}
	}

	if neg {
		d = -d
	}
	return d, nil
}

func (d JsonDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(durationToISO8601(time.Duration(d)))
}

func (d *JsonDuration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = JsonDuration(time.Duration(value))
		return nil
	case string:
		// Try ISO8601 first, then fall back to Go duration
		if tmp, err := parseISO8601(value); err == nil {
			*d = JsonDuration(tmp)
			return nil
		}
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = JsonDuration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}