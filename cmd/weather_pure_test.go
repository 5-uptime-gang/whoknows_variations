package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildOpenMeteoURL(t *testing.T) {
	cases := []struct {
		lat, lon, days, units string
		want                  []string
	}{
		{"55.6761", "12.5683", "5", "metric",
			[]string{"latitude=55.6761", "longitude=12.5683", "forecast_days=5", "temperature_unit=celsius"}},
		{"10", "20", "7", "imperial",
			[]string{"latitude=10", "longitude=20", "forecast_days=7", "temperature_unit=fahrenheit"}},
	}
	for _, c := range cases {
		u := buildOpenMeteoURL(c.lat, c.lon, c.days, c.units)
		if !strings.HasPrefix(u, "https://api.open-meteo.com/v1/forecast?") {
			t.Fatalf("bad base url: %s", u)
		}
		for _, want := range c.want {
			if !strings.Contains(u, want) {
				t.Fatalf("url %q missing %q", u, want)
			}
		}
	}
}

func TestNormalizeOpenMeteo(t *testing.T) {
	raw := []byte(`{
	  "daily": {
	    "time": ["2025-11-01", "2025-11-02"],
	    "weathercode": [3, 80],
	    "temperature_2m_max": [12.5, 11.0],
	    "temperature_2m_min": [6.2, 5.1]
	  },
	  "timezone": "Europe/Copenhagen"
	}`)
	out, tz, err := normalizeOpenMeteo(raw)
	if err != nil {
		t.Fatalf("normalize error: %v", err)
	}
	if tz != "Europe/Copenhagen" {
		t.Fatalf("tz=%s", tz)
	}
	if len(out.Days) != 2 {
		b, _ := json.Marshal(out)
		t.Fatalf("want 2 days, got %d: %s", len(out.Days), string(b))
	}
	d0 := out.Days[0]
	if d0.Date != "2025-11-01" || d0.TMax != 12.5 || d0.TMin != 6.2 || d0.Code != 3 {
		t.Fatalf("unexpected day0: %#v", d0)
	}
}
