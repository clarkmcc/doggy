package doggy

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	client             StatsdClient
	initOnce           sync.Once
	PanicOnStatsdError = true
)

// InitClient is a function that knows how to initialize a dogstatsd client. By default, it will
// lazily initialize a client when we record our first metric and attempts to connect to the address
// contained in the 'DOGSTATSD_ADDR' environment variable.
var InitClient = func() (*statsd.Client, error) {
	client, err := statsd.New(os.Getenv("DOGSTATSD_ADDR"))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func initClient() {
	initOnce.Do(func() {
		var err error
		client, err = InitClient()
		if err != nil {
			if PanicOnStatsdError {
				panic(err)
			}
			client = &mockClient{}
		}
	})
}

// Metric is a common group of fields across every type of metric (counter, gauge, histogram, etc)
type Metric struct {
	Namespace   string
	ServiceName string
	MetricName  string
	Tags        Tags
}

// getName returns the full name of a metric.
func (m Metric) getName() string {
	return strings.Trim(strings.Join([]string{
		m.Namespace, m.ServiceName, m.MetricName,
	}, "."), " .")
}

type CounterMetric struct {
	Metric
}

// Count tracks how many times something happened per second and panics if no statsd client exists
func (m CounterMetric) Count(value int, options ...MetricOption) {
	err := m.CountE(value, options...)
	if err != nil {
		panic(err)
	}
}

// CountE tracks how many times something happened per second.
func (m CounterMetric) CountE(value int, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Count(m.getName(), int64(value), opts.getTags(), opts.SampleRate)
}

type HistogramMetric struct {
	Metric
}

// Histogram tracks the statistical distribution of a set of values on each host and panics if no
// statsd client exists.
func (m HistogramMetric) Histogram(value float64, options ...MetricOption) {
	err := m.HistogramE(value, options...)
	if err != nil {
		panic(err)
	}
}

// HistogramE tracks the statistical distribution of a set of values on each host.
func (m HistogramMetric) HistogramE(value float64, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Histogram(m.getName(), value, opts.getTags(), opts.SampleRate)
}

type GaugeMetric struct {
	Metric
}

// Gauge measures the value of a metric at a particular time and panics if no statsd client exists.
func (m GaugeMetric) Gauge(value float64, options ...MetricOption) {
	err := m.GaugeE(value, options...)
	if err != nil {
		panic(err)
	}
}

// GaugeE measures the value of a metric at a particular time.
func (m GaugeMetric) GaugeE(value float64, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Gauge(m.getName(), value, opts.getTags(), opts.SampleRate)
}

type DistributionMetric struct {
	Metric
}

// Distribution tracks the statistical distribution of a set of values across your infrastructure
// and panics if no statsd client exists.
func (m DistributionMetric) Distribution(value float64, options ...MetricOption) {
	err := m.DistributionE(value, options...)
	if err != nil {
		panic(err)
	}
}

// DistributionE tracks the statistical distribution of a set of values across your infrastructure.
func (m DistributionMetric) DistributionE(value float64, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Distribution(m.getName(), value, opts.getTags(), opts.SampleRate)
}

type TimingMetric struct {
	Metric
}

// Timing sends timing information and panics if no statsd client exists.
func (m TimingMetric) Timing(value time.Duration, options ...MetricOption) {
	err := m.TimingE(value, options...)
	if err != nil {
		panic(err)
	}
}

// TimingE sends timing information.
func (m TimingMetric) TimingE(value time.Duration, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Timing(m.getName(), value, opts.getTags(), opts.SampleRate)
}

func NewMetric[T CounterMetric | HistogramMetric | GaugeMetric | DistributionMetric | TimingMetric](namespace, serviceName, metricName string, options ...MetricOption) (out T) {
	opts := buildMetricOptions(options...)
	switch any(out).(type) {
	case CounterMetric:
		return T(CounterMetric{Metric{
			Namespace:   namespace,
			ServiceName: serviceName,
			MetricName:  metricName,
			Tags:        opts.Tags,
		}})
	case HistogramMetric:
		return T(HistogramMetric{Metric{
			Namespace:   namespace,
			ServiceName: serviceName,
			MetricName:  metricName,
			Tags:        opts.Tags,
		}})
	case GaugeMetric:
		return T(GaugeMetric{Metric{
			Namespace:   namespace,
			ServiceName: serviceName,
			MetricName:  metricName,
			Tags:        opts.Tags,
		}})
	case DistributionMetric:
		return T(DistributionMetric{Metric{
			Namespace:   namespace,
			ServiceName: serviceName,
			MetricName:  metricName,
			Tags:        opts.Tags,
		}})
	case TimingMetric:
		return T(TimingMetric{Metric{
			Namespace:   namespace,
			ServiceName: serviceName,
			MetricName:  metricName,
			Tags:        opts.Tags,
		}})
	default:
		panic(fmt.Sprintf("unsupported type %T", out))
	}
}
