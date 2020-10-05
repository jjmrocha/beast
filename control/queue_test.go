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

	"github.com/jjmrocha/beast/client"
)

func TestNewFiFo(t *testing.T) {
	underTest := NewFifo()

	if underTest.len != 0 {
		t.Errorf("len must be 0")
	}

	if underTest.first != nil {
		t.Errorf("first must be nil")
	}

	if underTest.last != nil {
		t.Errorf("last must be nil")
	}
}

func TestPush(t *testing.T) {
	// given
	underTest := NewFifo()
	var tests = []struct {
		status        int
		expectedLen   int
		expectedFirst int
		expectedLast  int
	}{
		{1, 1, 1, 1},
		{2, 2, 1, 2},
		{3, 3, 1, 3},
	}
	// then
	for _, test := range tests {
		underTest.Push(buildResponse(test.status))

		if underTest.len != test.expectedLen {
			t.Errorf("len: got %v expected %v", underTest.len, test.expectedLen)
		}

		if underTest.first.Value.StatusCode != test.expectedFirst {
			t.Errorf("first: got %v expected %v", underTest.first.Value.StatusCode, test.expectedFirst)
		}

		if underTest.last.Value.StatusCode != test.expectedLast {
			t.Errorf("last: got %v expected %v", underTest.last.Value.StatusCode, test.expectedLast)
		}
	}
}

func TestPop(t *testing.T) {
	// given
	underTest := NewFifo()
	underTest.Push(buildResponse(1))
	underTest.Push(buildResponse(2))
	underTest.Push(buildResponse(3))
	// then

	// Remove first
	result := underTest.Pop()

	if result.Value.StatusCode != 1 {
		t.Errorf("status: got %v expected %v", result, 1)
	}

	if underTest.len != 2 {
		t.Errorf("len: got %v expected %v", underTest.len, 2)
	}

	if underTest.first.Value.StatusCode != 2 {
		t.Errorf("first: got %v expected %v", underTest.first.Value.StatusCode, 2)
	}

	if underTest.last.Value.StatusCode != 3 {
		t.Errorf("last: got %v expected %v", underTest.last.Value.StatusCode, 3)
	}

	// Remove second
	result = underTest.Pop()

	if result.Value.StatusCode != 2 {
		t.Errorf("status: got %v expected %v", result, 2)
	}

	if underTest.len != 1 {
		t.Errorf("len: got %v expected %v", underTest.len, 1)
	}

	if underTest.first.Value.StatusCode != 3 {
		t.Errorf("first: got %v expected %v", underTest.first.Value.StatusCode, 3)
	}

	if underTest.last.Value.StatusCode != 3 {
		t.Errorf("last: got %v expected %v", underTest.last.Value.StatusCode, 3)
	}

	// Remove last
	result = underTest.Pop()

	if result.Value.StatusCode != 3 {
		t.Errorf("status: got %v expected %v", result, 3)
	}

	if underTest.len != 0 {
		t.Errorf("len: got %v expected %v", underTest.len, 0)
	}

	if underTest.first != nil {
		t.Errorf("first: got %v expected %v", underTest.first, nil)
	}

	if underTest.last != nil {
		t.Errorf("last: got %v expected %v", underTest.last, nil)
	}

	// Nothing to return
	result = underTest.Pop()

	if result != nil {
		t.Errorf("status: got %v expected %v", result, nil)
	}

	if underTest.len != 0 {
		t.Errorf("len: got %v expected %v", underTest.len, 0)
	}

	if underTest.first != nil {
		t.Errorf("first: got %v expected %v", underTest.first, nil)
	}

	if underTest.last != nil {
		t.Errorf("last: got %v expected %v", underTest.last, nil)
	}
}

func TestLen(t *testing.T) {
	// given
	underTest := NewFifo()
	var tests = []struct {
		push     bool
		expected int
	}{
		{true, 1},
		{true, 2},
		{false, 1},
		{false, 0},
	}
	// then
	for _, test := range tests {
		if test.push {
			underTest.Push(buildResponse(1))
		} else {
			underTest.Pop()
		}

		result := underTest.Len()

		if result != test.expected {
			t.Errorf("got %v expected %v", result, test.expected)
		}
	}
}

func buildResponse(status int) *client.Response {
	return &client.Response{
		StatusCode: status,
	}
}
