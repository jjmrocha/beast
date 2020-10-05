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

package control

import (
	"testing"
)

func TestInit(t *testing.T) {
	// given
	var tests = []struct {
		input    bool
		expected bool
	}{
		{true, true},
		{false, false},
	}
	// then
	for _, test := range tests {
		underTest := NewAtomicBool(test.input)
		result := underTest.Get()
		if result != test.expected {
			t.Errorf("got %v expected %v", result, test.expected)
		}
	}
}

func TestSetGet(t *testing.T) {
	// given
	var tests = []struct {
		input    bool
		expected bool
	}{
		{true, true},
		{false, false},
	}
	underTest := NewAtomicBool(false)
	// then
	for _, test := range tests {
		underTest.Set(test.input)
		result := underTest.Get()
		if result != test.expected {
			t.Errorf("got %v expected %v", result, test.expected)
		}
	}
}
