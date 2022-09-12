package doggy

import (
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type MetricOptions struct {
	Tags       Tags
	SampleRate float64
}

func (o MetricOptions) getTags(merge ...Tags) (out []string) {
	return o.Tags.getTags(merge...)
}

func buildMetricOptions(options ...MetricOption) *MetricOptions {
	opts := MetricOptions{
		SampleRate: 1,
	}
	for _, opt := range options {
		opt.applyMetricOption(&opts)
	}
	return &opts
}

type MetricOption interface {
	applyMetricOption(*MetricOptions)
}

var _ MetricOption = &Tags{}
var _ ServiceCheckOption = &Tags{}

type Tags map[string]string

func (t Tags) applyServiceCheckOption(options *ServiceCheckOptions) {
	options.Tags = t
}

func (t Tags) applyMetricOption(options *MetricOptions) {
	options.Tags = t
}

// MergeInto merges t1 into t2 giving priority to any matching keys in t1
func (t1 Tags) MergeInto(t2 Tags) {
	for k, v := range t1 {
		t2[k] = v
	}
}

func (t1 Tags) getTags(t2 ...Tags) (out []string) {
	if len(t2) > 0 {
		t2[0].MergeInto(t1)
	}
	for k, v := range t1 {
		out = append(out, fmt.Sprintf("%v:%v", k, v))
	}
	return
}

type ServiceCheckOptions struct {
	Tags     Tags
	Hostname string
	Message  string
}

type ServiceCheckOption interface {
	applyServiceCheckOption(check *ServiceCheckOptions)
}

var _ ServiceCheckOption = Hostname("")

type Hostname string

func (h Hostname) applyServiceCheckOption(check *ServiceCheckOptions) {
	check.Hostname = string(h)
}

var _ ServiceCheckOption = Message("")

type Message string

func (m Message) applyServiceCheckOption(check *ServiceCheckOptions) {
	check.Message = string(m)
}

func buildServiceCheckOptions(options ...ServiceCheckOption) *ServiceCheckOptions {
	opts := ServiceCheckOptions{}
	for _, opt := range options {
		opt.applyServiceCheckOption(&opts)
	}
	return &opts
}

type TracerStartSpanOption interface {
	intoStartSpanOption() tracer.StartSpanOption
}

var _ TracerStartSpanOption = StartSpanOption(nil)

type StartSpanOption tracer.StartSpanOption

func (s StartSpanOption) intoStartSpanOption() tracer.StartSpanOption {
	return tracer.StartSpanOption(s)
}

var _ TracerStartSpanOption = TraceCache("")

type TraceCache string

func (t TraceCache) intoStartSpanOption() tracer.StartSpanOption {
	return tracer.SpanType("cache")
}

var _ TracerStartSpanOption = TraceDatabase("")

type TraceDatabase string

func (t TraceDatabase) intoStartSpanOption() tracer.StartSpanOption {
	return tracer.SpanType("db")
}
