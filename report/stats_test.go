/*
 * Copyright 2019 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
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
)

func TestPercentage(t *testing.T) {
	// given
	ds := durationSlice{
		time.Duration(0),
		time.Duration(1),
		time.Duration(2),
		time.Duration(3),
		time.Duration(4),
		time.Duration(5),
		time.Duration(6),
		time.Duration(7),
		time.Duration(8),
		time.Duration(9),
	}
	var tests = []struct {
		input    int
		expected time.Duration
	}{
		{20, time.Duration(1)},
		{40, time.Duration(3)},
		{60, time.Duration(5)},
		{80, time.Duration(7)},
	}
	// then
	for _, test := range tests {
		result := ds.percentage(test.input)
		if result != test.expected {
			t.Errorf("got %v expected %v for percentage %v", result, test.expected, test.input)
		}
	}
}

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
