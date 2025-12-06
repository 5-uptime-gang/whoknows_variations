package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type wxDay struct {
	Date string  `json:"date"`
	TMin float64 `json:"tMin"`
	TMax float64 `json:"tMax"`
	Code int     `json:"code"`
}

type wxOut struct {
	Source   string  `json:"source"`
	Updated  string  `json:"updated"`
	Timezone string  `json:"timezone"`
	Days     []wxDay `json:"daily"`
}

type cacheEntry struct {
	payload []byte
	expires time.Time
}

var (
	cacheTTL = 15 * time.Minute
	wxCache  = struct {
		mu sync.RWMutex
		m  map[string]cacheEntry
	}{m: make(map[string]cacheEntry)}
)

func apiWeather(c *gin.Context) {
	const (
		lat   = "55.6761" // Copenhagen
		lon   = "12.5683"
		days  = "5"
		units = "metric"
	)
	key := lat + "|" + lon + "|" + days + "|" + units

	if payload, ok := getCached(key); ok {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", payload)
		return
	}

	payload, ok := fetchWeather(key, lat, lon, days, units)
	if ok {
		c.Header("X-Cache", "MISS")
		c.Data(http.StatusOK, "application/json", payload)
		return
	}

	if payload, ok := getCached(key); ok {
		c.Header("X-Cache", "STALE")
		c.Data(http.StatusOK, "application/json", payload)
		return
	}

	payload, _ = json.Marshal(StandardResponse{Data: struct{}{}})
	c.Header("X-Cache", "EMPTY")
	c.Data(http.StatusOK, "application/json", payload)
}

func getCached(key string) ([]byte, bool) {
	wxCache.mu.RLock()
	defer wxCache.mu.RUnlock()
	ce, ok := wxCache.m[key]
	if !ok || time.Now().After(ce.expires) {
		return nil, false
	}
	return ce.payload, true
}

func fetchWeather(key, lat, lon, days, units string) ([]byte, bool) {
	url := buildOpenMeteoURL(lat, lon, days, units)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, false
	}
	req.Header.Set("User-Agent", "who-knows-weather/1.0")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, false
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("[WEATHER] Error closing response body: %v", cerr)
		}
	}()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false
	}

	out, tz, err := normalizeOpenMeteo(raw)
	if err != nil {
		return nil, false
	}

	out.Source = "open-meteo"
	out.Updated = time.Now().UTC().Format(time.RFC3339)
	out.Timezone = tz

	payload, err := json.Marshal(StandardResponse{Data: out})
	if err != nil {
		return nil, false
	}

	wxCache.mu.Lock()
	wxCache.m[key] = cacheEntry{payload, time.Now().Add(cacheTTL)}
	wxCache.mu.Unlock()

	return payload, true
}

func buildOpenMeteoURL(lat, lon, days, units string) string {
	tempUnit := "celsius"
	if units == "imperial" {
		tempUnit = "fahrenheit"
	}
	return "https://api.open-meteo.com/v1/forecast?latitude=" + lat +
		"&longitude=" + lon +
		"&timezone=auto&daily=weathercode,temperature_2m_max,temperature_2m_min" +
		"&forecast_days=" + days +
		"&temperature_unit=" + tempUnit
}

type omResp struct {
	Daily struct {
		Time        []string  `json:"time"`
		WeatherCode []int     `json:"weathercode"`
		TempMax     []float64 `json:"temperature_2m_max"`
		TempMin     []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
	Timezone string `json:"timezone"`
}

func normalizeOpenMeteo(raw []byte) (wxOut, string, error) {
	var r omResp
	if err := json.Unmarshal(raw, &r); err != nil {
		return wxOut{}, "", err
	}
	o := wxOut{}
	for i := range r.Daily.Time {
		o.Days = append(o.Days, wxDay{
			Date: r.Daily.Time[i],
			TMin: r.Daily.TempMin[i],
			TMax: r.Daily.TempMax[i],
			Code: r.Daily.WeatherCode[i],
		})
	}
	return o, r.Timezone, nil
}
