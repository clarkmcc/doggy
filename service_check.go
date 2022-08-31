package doggy

import (
	"github.com/DataDog/datadog-go/statsd"
	"strings"
	"time"
)

type ServiceCheckStatus byte

const (
	// Ok is the "ok" ServiceCheck status
	Ok ServiceCheckStatus = 0
	// Warn is the "warning" ServiceCheck status
	Warn ServiceCheckStatus = 1
	// Critical is the "critical" ServiceCheck status
	Critical ServiceCheckStatus = 2
	// Unknown is the "unknown" ServiceCheck status
	Unknown ServiceCheckStatus = 3
)

type ServiceCheck struct {
	Namespace   string
	ServiceName string
	CheckName   string
}

func (s ServiceCheck) getName() string {
	return strings.Trim(strings.Join([]string{
		s.Namespace, s.ServiceName, s.CheckName,
	}, "."), " .")
}

func (s ServiceCheck) UpdateStatusE(status ServiceCheckStatus, options ...ServiceCheckOption) error {
	initClient()
	opts := buildServiceCheckOptions(options...)
	return client.ServiceCheck(&statsd.ServiceCheck{
		Name:      s.getName(),
		Status:    statsd.ServiceCheckStatus(status),
		Timestamp: time.Now(),
		Hostname:  opts.Hostname,
		Message:   opts.Message,
		Tags:      opts.Tags.getTags(),
	})
}

// UpdateStatus updates the status of this service check
func (s ServiceCheck) UpdateStatus(status ServiceCheckStatus, options ...ServiceCheckOption) {
	err := s.UpdateStatusE(status, options...)
	if err != nil {
		panic(err)
	}
}

func NewServiceCheck(namespace, service, check string) *ServiceCheck {
	return &ServiceCheck{
		Namespace:   namespace,
		ServiceName: service,
		CheckName:   check,
	}
}
