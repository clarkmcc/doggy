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
		namespace, service, metric, expected string
	}{
		"happy path": {
			service:  "foo",
			metric:   "bar",
			expected: "foo.bar",
		},
		"happy path with prefix": {
			namespace: "ns",
			service:   "foo",
			metric:    "bar",
			expected:  "ns.foo.bar",
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			assert.Equal(t, test.expected, (Metric{
				Namespace:   test.namespace,
				ServiceName: test.service,
				MetricName:  test.metric,
			}).getName())
		})
	}
}
