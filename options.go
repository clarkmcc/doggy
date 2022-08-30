package doggy

import "fmt"

type MetricOptions struct {
	Tags       Tags
	SampleRate float64
}

func (o MetricOptions) getTags(merge ...Tags) (out []string) {
	tags := o.Tags
	if len(merge) > 0 {
		merge[0].MergeInto(tags)
	}
	for k, v := range tags {
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

func (t Tags) apply(options *MetricOptions) {
	options.Tags = t
}

// MergeFrom merges t1 into t2 giving priority to any matching keys in t1
func (t1 Tags) MergeInto(t2 Tags) {
	for k, v := range t1 {
		t2[k] = v
	}
}
