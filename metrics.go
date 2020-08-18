package logr

import "errors"

const (
	DefMetricsUpdateFreqMillis = 15000 // 15 seconds
)

// Counter is a simple metrics sink that can only increment a value.
// Implementations are external to Logr and provided via `MetricsCollector`.
type Counter interface {
	// Inc increments the counter by 1. Use Add to increment it by arbitrary non-negative values.
	Inc()
	// Add adds the given value to the counter. It panics if the value is < 0.
	Add(float64)
}

// Gauge is a simple metrics sink that can receive values and increase or decrease.
// Implementations are external to Logr and provided via `MetricsCollector`.
type Gauge interface {
	// Set sets the Gauge to an arbitrary value.
	Set(float64)
	// Add adds the given value to the Gauge. (The value can be negative, resulting in a decrease of the Gauge.)
	Add(float64)
	// Sub subtracts the given value from the Gauge. (The value can be negative, resulting in an increase of the Gauge.)
	Sub(float64)
}

// MetricsCollector provides a way for users of this Logr package to have metrics pushed
// in an efficient way to any backend, e.g. Prometheus.
// For each target added to Logr, the supplied MetricsCollector will provide a Gauge
// and Counters that will be called frequently as logging occurs.
type MetricsCollector interface {
	// QueueSizeGauge returns a Gauge that will be updated by the named target.
	QueueSizeGauge(target string) Gauge
	// LoggedCounter returns a Counter that will be incremented by the named target.
	LoggedCounter(target string) Counter
	// ErrorCounter returns a Counter that will be incremented by the named target.
	ErrorCounter(target string) Counter
	// DroppedCounter returns a Counter that will be incremented by the named target.
	DroppedCounter(target string) Counter
	// BlockedCounter returns a Counter that will be incremented by the named target.
	BlockedCounter(target string) Counter
}

// TargetWithMetrics is a target that provides metrics.
type TargetWithMetrics interface {
	Metrics(collector MetricsCollector, updateFreqMillis int64)
}

// SetMetricsCollector enables metrics collection by supplying a MetricsCollector.
// The MetricsCollector provides counters and gauges that are updated by log targets.
// This MUST be called before any log targets are added.
func (logr *Logr) SetMetricsCollector(collector MetricsCollector) error {
	if logr.HasTargets() {
		return errors.New("Logr.SetMetricsCollector must be called before any targets are added.")
	}
	if collector == nil {
		return errors.New("collector cannot be nil")
	}

	logr.metrics = collector
	logr.queueSizeGauge = collector.QueueSizeGauge("_logr")
	logr.loggedCounter = collector.LoggedCounter("_logr")
	logr.errorCounter = collector.ErrorCounter("_logr")
	return nil
}
