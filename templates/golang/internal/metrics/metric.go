package metrics

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	UserRepoCalls *prometheus.CounterVec
	// Add other metrics here, e.g., ProductRepoCalls, ApiLatency, etc.
}

func NewMetrics() *Metrics {
	serviceName := os.Getenv("SERVICE_NAME")
	return &Metrics{
		UserRepoCalls: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_user_repository_calls_total", serviceName),
				Help: "Total number of calls to the user repository.",
			},
			[]string{"method", "status"}, // Labels for the counter
		),
		// Initialize other metrics here...
	}
}
