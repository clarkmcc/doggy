package doggy

import "fmt"

type MetricOptions struct {
	Tags       map[string]string
	SampleRate float64
}

func (o MetricOptions) getTags() (out []string) {
	for k, v := range o.Tags {
		out = append(out, fmt.Sprintf("%v:%v", k, v))
	}
	return
}

func buildMetricOptions(options ...MetricOption) *MetricOptions {
	opts := MetricOptions{
		SampleRate: 1,
	}
	for _, opt := range options {
		opt.apply(&opts)
	}
	return &opts
}

type MetricOption interface {
	apply(*MetricOptions)
}

var _ MetricOption = &Tags{}

type Tags map[string]string

func (w Tags) apply(options *MetricOptions) {
	options.Tags = w
}
