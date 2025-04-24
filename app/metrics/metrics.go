package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	counter *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

func New() *Metrics {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "counter",
		Help: "Count status",
	}, []string{"status"})

	latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "latency_seconds",
		Help:    "Duration of the latency in seconds",
		Buckets: prometheus.ExponentialBuckets(0.05, 1.3, 15),
	}, []string{"name"})

	prometheus.MustRegister(counter)
	prometheus.MustRegister(latency)

	return &Metrics{
		latency: latency,
		counter: counter,
	}
}

func (m *Metrics) Latency(labels string) *prometheus.Timer {
	return prometheus.NewTimer(m.latency.WithLabelValues(labels))
}

func (m *Metrics) Count(status string) {
	m.counter.WithLabelValues(status).Inc()
}
