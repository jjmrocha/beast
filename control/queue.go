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

import "github.com/jjmrocha/beast/client"

// Element is a value on the queue
type Element struct {
	next  *Element
	Value *client.Response
}

// Fifo is a linked list
type Fifo struct {
	first *Element
	last  *Element
	len   int
}

// NewFifo creates a new queue
func NewFifo() *Fifo {
	return &Fifo{
		first: nil,
		last:  nil,
		len:   0,
	}
}

// Len returns the size of the queue
func (q *Fifo) Len() int {
	return q.len
}

// Push adds a element to the end of the queue
func (q *Fifo) Push(value *client.Response) {
	element := &Element{
		next:  nil,
		Value: value,
	}

	if q.len == 0 {
		q.first = element
	} else {
		q.last.next = element
	}

	q.last = element
	q.len = q.len + 1
}

// Pop retrives the fisrt element of the queue
func (q *Fifo) Pop() *Element {
	if q.len == 0 {
		return nil
	}

	element := q.first
	q.first = element.next
	q.len = q.len - 1

	if q.len == 0 {
		q.last = nil
	}

	return element
}
