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

	"github.com/stretchr/testify/assert"
)

func TestNewAtomicBool(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected bool
	}{
		{name: "starts true", input: true, expected: true},
		{name: "starts false", input: false, expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			input := tc.input
			// when
			result := NewAtomicBool(input)
			// then
			assert.Equal(t, tc.expected, result.Get())
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name     string
		initial  bool
		input    bool
		expected bool
	}{
		{name: "false to true", initial: false, input: true, expected: true},
		{name: "true to false", initial: true, input: false, expected: false},
		{name: "true stays true", initial: true, input: true, expected: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			atomicBool := NewAtomicBool(tc.initial)
			// when
			atomicBool.Set(tc.input)
			// then
			assert.Equal(t, tc.expected, atomicBool.Get())
		})
	}
}
