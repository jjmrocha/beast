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
)

func TestSuccess(t *testing.T) {
	// given
	var tests = []struct {
		input    int
		expected bool
	}{
		{100, false},
		{122, false},
		{200, true},
		{201, true},
		{300, false},
		{307, false},
		{400, false},
		{404, false},
		{429, false},
		{500, false},
		{502, false},
	}
	// then
	for _, test := range tests {
		result := success(test.input)
		if result != test.expected {
			t.Errorf("got %v expected %v for status %v", result, test.expected, test.input)
		}
	}
}
