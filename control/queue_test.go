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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFifo(t *testing.T) {
	// when
	result := NewFifo()
	// then
	assert.Equal(t, 0, result.Len())
	assert.Nil(t, result.Pop())
}

func TestPush(t *testing.T) {
	t.Run("adds to an empty queue", func(t *testing.T) {
		// given
		q := NewFifo()
		// when
		q.Push(response(1))
		// then
		assert.Equal(t, 1, q.Len())
	})

	t.Run("adds at the end", func(t *testing.T) {
		// given
		q := NewFifo()
		q.Push(response(1))
		q.Push(response(2))
		// when
		q.Push(response(3))
		// then
		assert.Equal(t, 3, q.Len())
		assert.Equal(t, []int{1, 2, 3}, drain(q))
	})
}

func TestPop(t *testing.T) {
	t.Run("returns the first pushed element", func(t *testing.T) {
		// given
		q := NewFifo()
		q.Push(response(1))
		q.Push(response(2))
		q.Push(response(3))
		// when
		result := q.Pop()
		// then
		require.NotNil(t, result)
		assert.Equal(t, 1, result.Value.StatusCode)
		assert.Equal(t, 2, q.Len())
	})

	t.Run("returns the last element and leaves the queue empty", func(t *testing.T) {
		// given
		q := NewFifo()
		q.Push(response(1))
		// when
		result := q.Pop()
		// then
		require.NotNil(t, result)
		assert.Equal(t, 1, result.Value.StatusCode)
		assert.Equal(t, 0, q.Len())
		assert.Nil(t, q.Pop())
	})

	t.Run("returns nil on an empty queue", func(t *testing.T) {
		// given
		q := NewFifo()
		// when
		result := q.Pop()
		// then
		assert.Nil(t, result)
		assert.Equal(t, 0, q.Len())
	})

	t.Run("works again after the queue was emptied", func(t *testing.T) {
		// given
		q := NewFifo()
		q.Push(response(1))
		q.Pop()
		q.Push(response(2))
		// when
		result := q.Pop()
		// then
		require.NotNil(t, result)
		assert.Equal(t, 2, result.Value.StatusCode)
	})
}

func TestLen(t *testing.T) {
	tests := []struct {
		name     string
		pushes   int
		pops     int
		expected int
	}{
		{name: "new queue", pushes: 0, pops: 0, expected: 0},
		{name: "after two pushes", pushes: 2, pops: 0, expected: 2},
		{name: "after two pushes and one pop", pushes: 2, pops: 1, expected: 1},
		{name: "after two pushes and two pops", pushes: 2, pops: 2, expected: 0},
		{name: "pop on empty does not go below zero", pushes: 0, pops: 1, expected: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			q := NewFifo()
			for i := range tc.pushes {
				q.Push(response(i))
			}
			for range tc.pops {
				q.Pop()
			}
			// when
			result := q.Len()
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func response(status int) *client.Response {
	return &client.Response{StatusCode: status}
}

func drain(q *Fifo) []int {
	var result []int
	for element := q.Pop(); element != nil; element = q.Pop() {
		result = append(result, element.Value.StatusCode)
	}
	return result
}
