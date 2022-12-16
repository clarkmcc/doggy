package doggy

import (
	"github.com/DataDog/datadog-go/statsd"
	"time"
)

type StatsdClient interface {
	Count(name string, value int64, tags []string, rate float64) error
	Histogram(name string, value float64, tags []string, rate float64) error
	Gauge(name string, value float64, tags []string, rate float64) error
	Distribution(name string, value float64, tags []string, rate float64) error
	Timing(name string, value time.Duration, tags []string, rate float64) error
	ServiceCheck(sc *statsd.ServiceCheck) error
	Event(event *statsd.Event) error
}

var _ StatsdClient = &statsd.Client{}
var _ StatsdClient = &mockClient{}

type mockClient struct{}

func (m *mockClient) Count(name string, value int64, tags []string, rate float64) error {
	return nil
}

func (m *mockClient) Histogram(name string, value float64, tags []string, rate float64) error {
	return nil
}

func (m *mockClient) Gauge(name string, value float64, tags []string, rate float64) error {
	return nil
}

func (m *mockClient) Distribution(name string, value float64, tags []string, rate float64) error {
	return nil
}

func (m *mockClient) Timing(name string, value time.Duration, tags []string, rate float64) error {
	return nil
}

func (m *mockClient) ServiceCheck(sc *statsd.ServiceCheck) error {
	return nil
}

func (m *mockClient) Event(event *statsd.Event) error {
	return nil
}
