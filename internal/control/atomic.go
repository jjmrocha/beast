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

import "sync/atomic"

// AtomicBool is an atomic Boolean
type AtomicBool int32

const trueValue = 1
const falseValue = 0

// NewAtomicBool creates a new AtomicBool
func NewAtomicBool(value bool) *AtomicBool {
	atomicBool := new(AtomicBool)
	atomicBool.Set(value)
	return atomicBool
}

// Set can be use to change the value of a AtomicBool
func (atomicBool *AtomicBool) Set(value bool) {
	var atomicValue int32 = falseValue

	if value {
		atomicValue = trueValue
	}

	atomic.StoreInt32((*int32)(atomicBool), atomicValue)
}

// Get reads a value of a AtomicBool
func (atomicBool *AtomicBool) Get() bool {
	atomicValue := atomic.LoadInt32((*int32)(atomicBool))
	return atomicValue == trueValue
}
