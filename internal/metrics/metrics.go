package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Serve exposes metrics via http on the provided addr
func Serve(addr string) error {

	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}

	return nil
}

// Counters returns a map of counters from given names
func Counters(namespace, subsystem string, names []string) map[string]prometheus.Counter {

	counters := make(map[string]prometheus.Counter)

	for _, name := range names {
		counters[name] = promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
		})
	}

	return counters
}
