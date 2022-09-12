package doggy

import (
	"context"
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"os"
)

// StartTracer starts the tracer with the correct env variables
func StartTracer(serviceName, serviceVersion string, options ...tracer.StartOption) error {
	var opts []tracer.StartOption
	if len(serviceName) == 0 {
		return fmt.Errorf("service name required by tracer")
	}
	opts = append(opts, tracer.WithService(serviceName))
	if addr := os.Getenv("DOGSTATSD_APM_ADDR"); len(addr) > 0 {
		opts = append(opts, tracer.WithAgentAddr(addr))
	}
	if len(serviceVersion) > 0 {
		opts = append(opts, tracer.WithServiceVersion(serviceVersion))
	}
	if env := os.Getenv("ENV"); len(env) > 0 {
		opts = append(opts, tracer.WithEnv(env))
	}
	tracer.Start(append(opts, options...)...)
	return nil
}

// StopTracer stops the started tracer. Subsequent calls are valid but become no-op.
func StopTracer() {
	tracer.Stop()
}

// StartSpanFromContext starts a new span which is a child of an existing span if one has been injected
func StartSpanFromContext(ctx context.Context, operationName string, opts ...tracer.StartSpanOption) tracer.Span {
	var s tracer.Span
	s, ok := tracer.SpanFromContext(ctx)
	if !ok {
		s = tracer.StartSpan(operationName, append(opts, tracer.ChildOf(s.Context()))...)
	}
	return s
}
