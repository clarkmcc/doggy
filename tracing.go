package doggy

import (
	"context"
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"os"
	"runtime"
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
func StartSpanFromContext(ctx context.Context, operationName string, opts ...TracerStartSpanOption) tracer.Span {
	options := make([]tracer.StartSpanOption, len(opts))
	for i := 0; i < len(opts); i++ {
		options[i] = opts[i].intoStartSpanOption()
	}
	if len(operationName) == 0 {
		file, line, function := getCurrentFunc()
		operationName = function
		options = append(options, tracer.Tag("file", fmt.Sprintf("%v:%v", file, line)))
	}
	var span tracer.Span
	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		span = tracer.StartSpan(operationName, options...)
	} else {
		span = tracer.StartSpan(operationName, append(options, tracer.ChildOf(span.Context()))...)
	}
	return span
}

func getCurrentFunc() (string, int, string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.File, frame.Line, frame.Function
}
