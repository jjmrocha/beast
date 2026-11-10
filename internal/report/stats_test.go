/*
 * Copyright 2026 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package report

import (
	"bytes"
	"testing"
	"time"

	"github.com/jjmrocha/beast/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestTps(t *testing.T) {
	tests := []struct {
		name     string
		requests int
		duration time.Duration
		expected float64
	}{
		{name: "results divided by elapsed time", requests: 100, duration: 0, expected: 50},
		{name: "response times do not change the rate", requests: 1000, duration: 50 * time.Second, expected: 500},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			s := &Stats{
				requests:       tc.requests,
				duration:       tc.duration,
				executionStart: time.Now().Add(-2 * time.Second),
			}
			// when
			result := s.tps()
			// then
			assert.InDelta(t, tc.expected, result, tc.expected/100)
		})
	}
}

func TestPrintStats(t *testing.T) {
	t.Run("prints the number of executed requests", func(t *testing.T) {
		// given
		var printed bytes.Buffer
		s := statsOf(success(1), success(2))
		// when
		s.PrintStats(&printed)
		// then
		assert.Contains(t, printed.String(), "===== Stats =====\nExecuted requests: 2\n")
	})

	t.Run("prints the distribution of a success status with at least 5 results", func(t *testing.T) {
		// given
		var printed bytes.Buffer
		s := statsOf(success(5), success(1), success(4), success(2), success(3))
		// when
		s.PrintStats(&printed)
		// then
		result := printed.String()
		assert.Contains(t, result, "===== Status 200 =====\n5 requests, with avg response time of 3ms\n")
		assert.Contains(t, result, "- The fastest request took 1ms\n")
		assert.Contains(t, result, "- 80% of requests under 4ms\n")
		assert.Contains(t, result, "- The slowest request took 5ms\n")
	})

	t.Run("prints no distribution for fewer than 5 results", func(t *testing.T) {
		// given
		var printed bytes.Buffer
		s := statsOf(success(1), success(2), success(3), success(4))
		// when
		s.PrintStats(&printed)
		// then
		assert.Contains(t, printed.String(), "4 requests, with avg response time of 2.5ms\n")
		assert.NotContains(t, printed.String(), "And the following distribution")
	})

	t.Run("prints non-success statuses", func(t *testing.T) {
		// given
		var printed bytes.Buffer
		s := statsOf(&client.Response{StatusCode: 404, Duration: time.Millisecond})
		// when
		s.PrintStats(&printed)
		// then
		assert.Contains(t, printed.String(), "===== Non Success Status =====\nStatus 404: 1 requests\n")
	})

	t.Run("prints client errors", func(t *testing.T) {
		// given
		var printed bytes.Buffer
		s := statsOf(&client.Response{StatusCode: client.StatusTimeout, Duration: time.Second})
		// when
		s.PrintStats(&printed)
		// then
		assert.Contains(t, printed.String(), "===== Errors =====\n- Request timeout: 1 errors\n")
	})
}

type noProgress struct{}

func (noProgress) Update() {}

func statsOf(responses ...*client.Response) *Stats {
	s := NewStats(noProgress{})
	for _, response := range responses {
		s.Update(response)
	}
	return s
}

func success(milliseconds int) *client.Response {
	return &client.Response{StatusCode: 200, Duration: time.Duration(milliseconds) * time.Millisecond}
}
