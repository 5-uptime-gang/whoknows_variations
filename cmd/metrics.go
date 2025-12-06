package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type prometheusLabels struct {
	method string
	path   string
	status string
}

type histogram struct {
	counts []uint64
	sum    float64
	count  uint64
}

type metricsRecorder struct {
	mu      sync.Mutex
	buckets []float64
	counter map[prometheusLabels]uint64
	timing  map[prometheusLabels]*histogram
}

var defaultMetrics = newMetricsRecorder()

func newMetricsRecorder() *metricsRecorder {
	return &metricsRecorder{
		buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		counter: make(map[prometheusLabels]uint64),
		timing:  make(map[prometheusLabels]*histogram),
	}
}

func observeRequest(duration time.Duration, labels prometheusLabels) {
	defaultMetrics.observe(duration, labels)
}

func (m *metricsRecorder) observe(duration time.Duration, labels prometheusLabels) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter[labels]++

	hist, ok := m.timing[labels]
	if !ok {
		hist = &histogram{counts: make([]uint64, len(m.buckets))}
		m.timing[labels] = hist
	}

	seconds := duration.Seconds()
	hist.sum += seconds
	hist.count++
	for i, bound := range m.buckets {
		if seconds <= bound {
			hist.counts[i]++
		}
	}
}

func metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprint(w, defaultMetrics.render())
}

func (m *metricsRecorder) render() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var builder strings.Builder

	builder.WriteString("# HELP http_requests_total Total number of HTTP requests\n")
	builder.WriteString("# TYPE http_requests_total counter\n")
	for _, labels := range m.sortedCounterKeys() {
		builder.WriteString(fmt.Sprintf("http_requests_total%s %d\n", formatCounterLabels(labels), m.counter[labels]))
	}

	builder.WriteString("# HELP http_request_duration_seconds Time taken to serve HTTP requests\n")
	builder.WriteString("# TYPE http_request_duration_seconds histogram\n")
	for _, labels := range m.sortedTimingKeys() {
		hist := m.timing[labels]
		for i, bound := range m.buckets {
			builder.WriteString(fmt.Sprintf(
				"http_request_duration_seconds_bucket%s %d\n",
				formatHistogramLabels(labels, strconv.FormatFloat(bound, 'g', -1, 64)),
				hist.counts[i],
			))
		}
		builder.WriteString(fmt.Sprintf(
			"http_request_duration_seconds_bucket%s %d\n",
			formatHistogramLabels(labels, "+Inf"),
			hist.count,
		))
		builder.WriteString(fmt.Sprintf("http_request_duration_seconds_sum%s %.6f\n", formatCounterLabels(labels), hist.sum))
		builder.WriteString(fmt.Sprintf("http_request_duration_seconds_count%s %d\n", formatCounterLabels(labels), hist.count))
	}

	return builder.String()
}

func (m *metricsRecorder) sortedCounterKeys() []prometheusLabels {
	keys := make([]prometheusLabels, 0, len(m.counter))
	for k := range m.counter {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return labelKey(keys[i]) < labelKey(keys[j])
	})
	return keys
}

func (m *metricsRecorder) sortedTimingKeys() []prometheusLabels {
	keys := make([]prometheusLabels, 0, len(m.timing))
	for k := range m.timing {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return labelKey(keys[i]) < labelKey(keys[j])
	})
	return keys
}

func labelKey(label prometheusLabels) string {
	return label.method + "|" + label.path + "|" + label.status
}

func formatCounterLabels(labels prometheusLabels) string {
	return fmt.Sprintf("{method=\"%s\",path=\"%s\",status=\"%s\"}", escapeLabel(labels.method), escapeLabel(labels.path), escapeLabel(labels.status))
}

func formatHistogramLabels(labels prometheusLabels, bucket string) string {
	return fmt.Sprintf("{le=\"%s\",method=\"%s\",path=\"%s\",status=\"%s\"}", escapeLabel(bucket), escapeLabel(labels.method), escapeLabel(labels.path), escapeLabel(labels.status))
}

func escapeLabel(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}
