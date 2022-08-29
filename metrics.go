package doggy

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"sync"
	"time"
)

var (
	Prefix = ""

	client   *statsd.Client
	initOnce sync.Once
)

var InitClient = func() (*statsd.Client, error) {
	client, err := statsd.New("127.0.0.1:8125")
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
			panic(err)
		}
	})
}

type Metric struct {
	ServiceName string
	MetricName  string
}

func (m Metric) getName() string {
	name := m.ServiceName + "/" + m.MetricName
	if len(Prefix) > 0 {
		name = Prefix + name
	}
	return name
}

type CounterMetric struct {
	Metric
}

func (m CounterMetric) Count(value int, options ...MetricOption) {
	err := m.CountE(value, options...)
	if err != nil {
		panic(err)
	}
}

func (m CounterMetric) CountE(value int, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Count(m.getName(), int64(value), opts.getTags(), opts.SampleRate)
}

type HistogramMetric struct {
	Metric
}

func (m HistogramMetric) Histogram(value float64, options ...MetricOption) {
	err := m.HistogramE(value, options...)
	if err != nil {
		panic(err)
	}
}

func (m HistogramMetric) HistogramE(value float64, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Histogram(m.getName(), value, opts.getTags(), opts.SampleRate)
}

type GaugeMetric struct {
	Metric
}

func (m GaugeMetric) Gauge(value float64, options ...MetricOption) {
	err := m.GaugeE(value, options...)
	if err != nil {
		panic(err)
	}
}

func (m GaugeMetric) GaugeE(value float64, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Distribution(m.getName(), value, opts.getTags(), opts.SampleRate)
}

type DistributionMetric struct {
	Metric
}

func (m DistributionMetric) Distribution(value float64, options ...MetricOption) {
	err := m.DistributionE(value, options...)
	if err != nil {
		panic(err)
	}
}

func (m DistributionMetric) DistributionE(value float64, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Distribution(m.getName(), value, opts.getTags(), opts.SampleRate)
}

type TimingMetric struct {
	Metric
}

func (m TimingMetric) Timing(value time.Duration, options ...MetricOption) {
	err := m.TimingE(value, options...)
	if err != nil {
		panic(err)
	}
}

func (m TimingMetric) TimingE(value time.Duration, options ...MetricOption) error {
	initClient()
	opts := buildMetricOptions(options...)
	return client.Timing(m.getName(), value, opts.getTags(), opts.SampleRate)
}

func NewMetric[T CounterMetric | HistogramMetric | GaugeMetric | DistributionMetric](serviceName, metricName string) (out T) {
	switch any(out).(type) {
	case CounterMetric:
		return T(CounterMetric{Metric{
			ServiceName: serviceName,
			MetricName:  metricName,
		}})
	case HistogramMetric:
		return T(HistogramMetric{Metric{
			ServiceName: serviceName,
			MetricName:  metricName,
		}})
	case GaugeMetric:
		return T(GaugeMetric{Metric{
			ServiceName: serviceName,
			MetricName:  metricName,
		}})
	case DistributionMetric:
		return T(DistributionMetric{Metric{
			ServiceName: serviceName,
			MetricName:  metricName,
		}})
	case TimingMetric:
		return T(TimingMetric{Metric{
			ServiceName: serviceName,
			MetricName:  metricName,
		}})
	default:
		panic(fmt.Sprintf("unsupported type %T", out))
	}
}
