/*
 * Copyright 2019-20 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
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
	"testing"
	"time"

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
