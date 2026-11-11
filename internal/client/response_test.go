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

package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{name: "100 continue", input: 100, expected: false},
		{name: "122", input: 122, expected: false},
		{name: "199 just below 2xx", input: 199, expected: false},
		{name: "200 ok", input: 200, expected: true},
		{name: "201 created", input: 201, expected: true},
		{name: "299 last 2xx", input: 299, expected: true},
		{name: "300 multiple choices", input: 300, expected: false},
		{name: "307 temporary redirect", input: 307, expected: false},
		{name: "400 bad request", input: 400, expected: false},
		{name: "404 not found", input: 404, expected: false},
		{name: "429 too many requests", input: 429, expected: false},
		{name: "500 internal error", input: 500, expected: false},
		{name: "502 bad gateway", input: 502, expected: false},
		{name: "-400 client timeout", input: StatusTimeout, expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			response := &Response{StatusCode: tc.input}
			// when
			result := response.IsSuccess()
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFailed(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{name: "-500 unexpected error", input: StatusUnexpectedError, expected: true},
		{name: "-1 smallest negative", input: -1, expected: true},
		{name: "0 no status", input: 0, expected: false},
		{name: "200 ok", input: 200, expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			response := &Response{StatusCode: tc.input}
			// when
			result := response.Failed()
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFailureReason(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{name: "-100 generation error", input: StatusGenerationError, expected: "Request generation error"},
		{name: "-400 timeout", input: StatusTimeout, expected: "Request timeout"},
		{name: "-500 unexpected error", input: StatusUnexpectedError, expected: "Unexpected error"},
		{name: "-1 unknown client error", input: -1, expected: ""},
		{name: "200 not a client error", input: 200, expected: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			response := &Response{StatusCode: tc.input}
			// when
			result := response.FailureReason()
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}
