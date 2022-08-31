package doggy

import "github.com/DataDog/datadog-go/statsd"

// Event sends a new event to Datadog
func Event(event *statsd.Event) {
	initClient()
	err := client.Event(event)
	if err != nil {
		panic(err)
	}
}
