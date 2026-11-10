/*
 * Copyright 2020 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPercentage(t *testing.T) {
	ds := durationSlice{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	tests := []struct {
		input    int
		expected time.Duration
	}{
		{20, 1},
		{40, 3},
		{60, 5},
		{80, 7},
		{90, 8},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d percent of ten", tc.input), func(t *testing.T) {
			// given
			percentage := tc.input
			// when
			result := ds.percentage(percentage)
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}
