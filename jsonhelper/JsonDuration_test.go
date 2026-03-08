package jsonhelper

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDurationToISO8601(t *testing.T) {
	tests := []struct {
		dur  time.Duration
		want string
	}{
		{0, "PT0S"},
		{5 * time.Second, "PT5S"},
		{90 * time.Second, "PT1M30S"},
		{2 * time.Hour, "PT2H"},
		{2*time.Hour + 30*time.Minute + 15*time.Second, "PT2H30M15S"},
		{25 * time.Hour, "P1DT1H"},
		{48 * time.Hour, "P2D"},
		{49*time.Hour + 5*time.Minute, "P2DT1H5M"},
		{-3 * time.Hour, "-PT3H"},
	}
	for _, tt := range tests {
		got := durationToISO8601(tt.dur)
		if got != tt.want {
			t.Errorf("durationToISO8601(%v) = %q, want %q", tt.dur, got, tt.want)
		}
	}
}

func TestParseISO8601(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
	}{
		{"PT0S", 0},
		{"PT5S", 5 * time.Second},
		{"PT1M30S", 90 * time.Second},
		{"PT2H", 2 * time.Hour},
		{"PT2H30M15S", 2*time.Hour + 30*time.Minute + 15*time.Second},
		{"P1DT1H", 25 * time.Hour},
		{"P2D", 48 * time.Hour},
		{"P2DT1H5M", 49*time.Hour + 5*time.Minute},
		{"-PT3H", -3 * time.Hour},
	}
	for _, tt := range tests {
		got, err := parseISO8601(tt.input)
		if err != nil {
			t.Errorf("parseISO8601(%q) error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("parseISO8601(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseISO8601_Invalid(t *testing.T) {
	invalid := []string{"", "P", "T5S", "5S", "PXY"}
	for _, s := range invalid {
		if _, err := parseISO8601(s); err == nil {
			t.Errorf("parseISO8601(%q) expected error, got nil", s)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	d := JsonDuration(2*time.Hour + 30*time.Minute)
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"PT2H30M"` {
		t.Errorf("MarshalJSON = %s, want %q", b, "PT2H30M")
	}
}

func TestUnmarshalJSON_ISO8601(t *testing.T) {
	var d JsonDuration
	if err := json.Unmarshal([]byte(`"PT1H15M"`), &d); err != nil {
		t.Fatal(err)
	}
	want := JsonDuration(time.Hour + 15*time.Minute)
	if d != want {
		t.Errorf("UnmarshalJSON ISO8601 = %v, want %v", time.Duration(d), time.Duration(want))
	}
}

func TestUnmarshalJSON_GoDuration(t *testing.T) {
	var d JsonDuration
	if err := json.Unmarshal([]byte(`"1h15m"`), &d); err != nil {
		t.Fatal(err)
	}
	want := JsonDuration(time.Hour + 15*time.Minute)
	if d != want {
		t.Errorf("UnmarshalJSON Go duration = %v, want %v", time.Duration(d), time.Duration(want))
	}
}

func TestUnmarshalJSON_Number(t *testing.T) {
	var d JsonDuration
	ns := float64(5 * time.Second)
	b, _ := json.Marshal(ns)
	if err := json.Unmarshal(b, &d); err != nil {
		t.Fatal(err)
	}
	if time.Duration(d) != 5*time.Second {
		t.Errorf("UnmarshalJSON number = %v, want 5s", time.Duration(d))
	}
}

func TestRoundTrip(t *testing.T) {
	original := JsonDuration(3*time.Hour + 45*time.Minute + 10*time.Second)
	b, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded JsonDuration
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if original != decoded {
		t.Errorf("round trip: got %v, want %v", time.Duration(decoded), time.Duration(original))
	}
}
