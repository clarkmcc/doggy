// Copyright (C) 2022 Print Tracker, LLC - All Rights Reserved
//
// Unauthorized copying of this file, via any medium is strictly prohibited
// as this source code is proprietary and confidential. Dissemination of this
// information or reproduction of this material is strictly forbidden unless
// prior written permission is obtained from Print Tracker, LLC.

package doggy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetric_GetName(t *testing.T) {
	tests := map[string]struct {
		prefix, service, metric, expected string
	}{
		"happy path": {
			service:  "foo",
			metric:   "bar",
			expected: "foo.bar",
		},
		"happy path with prefix": {
			prefix:   "ptkr_io",
			service:  "foo",
			metric:   "bar",
			expected: "ptkr_io.foo.bar",
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			Prefix = test.prefix
			assert.Equal(t, test.expected, (Metric{
				ServiceName: test.service,
				MetricName:  test.metric,
			}).getName())
		})
	}
}
