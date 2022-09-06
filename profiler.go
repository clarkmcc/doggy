package doggy

import (
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"os"
)

// StartProfiler starts the profiler with the provided parameters. The service name is required,
// but the service version and environment are not required.
func StartProfiler(serviceName, serviceVersion string, options ...profiler.Option) error {
	// Maintain our own opts list so that we can put all our default options in front of
	// anything provided by the options parameter. This will ensure that anything provided
	// in the options parameter will override the manually set service name, version, or env
	var opts []profiler.Option
	if len(serviceName) == 0 {
		return fmt.Errorf("service name required by profiler")
	}
	opts = append(opts, profiler.WithService(serviceName))
	if addr := os.Getenv("DOGSTATSD_APM_ADDR"); len(addr) > 0 {
		opts = append(opts, profiler.WithAgentAddr(addr))
	}
	if len(serviceVersion) > 0 {
		opts = append(opts, profiler.WithVersion(serviceVersion))
	}
	if env := os.Getenv("ENV"); len(env) > 0 {
		opts = append(opts, profiler.WithEnv(env))
	}
	return profiler.Start(append(opts, options...)...)
}

// StopProfiler stops the profiler
func StopProfiler() {
	profiler.Stop()
}
