package doggy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTags(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		opts := buildMetricOptions(Tags{"foo": "bar"})
		assert.Equal(t, []string{"foo:bar"}, opts.getTags())
	})
	t.Run("merge", func(t *testing.T) {
		opts := buildMetricOptions(Tags{"foo": "bar"})
		assert.Equal(t, []string{
			"foo:baz",
			"bar:baz",
		}, opts.getTags(Tags{
			"foo": "baz",
			"bar": "baz",
		}))
	})
}
