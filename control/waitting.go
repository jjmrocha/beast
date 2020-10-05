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

// WaitCompletion allows a goroutine to request and wait for completion
type WaitCompletion struct {
	requested       *AtomicBool
	waittingChannel chan bool
}

// NewWaitCompletion creates a new WaitCompletion
func NewWaitCompletion() *WaitCompletion {
	return &WaitCompletion{
		requested:       NewAtomicBool(false),
		waittingChannel: make(chan bool),
	}
}

// Request and waits for completion
func (wc *WaitCompletion) Request() {
	wc.requested.Set(true)
	<-wc.waittingChannel
}

// WasRequested allows a goroutine if a request was send
func (wc *WaitCompletion) WasRequested() bool {
	return wc.requested.Get()
}

// Completed informs that the resquest was completed
func (wc *WaitCompletion) Completed() {
	wc.waittingChannel <- true
}
