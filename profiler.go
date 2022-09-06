package doggy

import (
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"os"
)

// StartProfiler starts the profiler with the provided parameters. The service name is required,
// but the service version and environment are not required.
func StartProfiler(serviceName, serviceVersion, env string, options ...profiler.Option) error {
	if len(serviceName) == 0 {
		return fmt.Errorf("service name required by profiler")
	}
	if addr := os.Getenv("DOGSTATSD_APM_ADDR"); len(addr) > 0 {
		options = append(options, profiler.WithAgentAddr(addr))
	}
	if len(serviceVersion) > 0 {
		options = append(options, profiler.WithVersion(serviceVersion))
	}
	if len(env) > 0 {
		options = append(options, profiler.WithEnv(env))
	} else if env = os.Getenv("ENV"); len(env) > 0 {
		options = append(options, profiler.WithEnv(env))
	}
	return profiler.Start(options...)
}

// StopProfiler stops the profiler
func StopProfiler() {
	profiler.Stop()
}
