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

// Package control provides functions to manage the execution of multiple goroutines
// in simultaneous
package control

import (
	"sync"

	"github.com/jjmrocha/beast/client"
)

type semaphore chan bool

func newSemaphore(size int) semaphore {
	return make(chan bool, size)
}

func (s semaphore) acquire() {
	s <- true
}

func (s semaphore) release() {
	<-s
}

// BControl is used to control the execution of multiple goroutines
type BControl struct {
	wg         sync.WaitGroup
	semaphore  semaphore
	outputChan chan *client.BResponse
}

// New creates a BControle
func New(nRequests, nParallel int) *BControl {
	ctrl := &BControl{
		semaphore:  newSemaphore(nParallel),
		outputChan: make(chan *client.BResponse, nRequests),
	}
	ctrl.wg.Add(nRequests)

	return ctrl
}

// Push sends the BResponse to the client goroutine
func (c *BControl) Push(response *client.BResponse) {
	c.outputChan <- response
}

// CloseWhenDone closes the OutputChannel when all goroutines finish execution
func (c *BControl) CloseWhenDone() {
	c.wg.Wait()
	close(c.outputChan)
}

// OutputChannel returns a channel for receiving the BResponse sent using Push
func (c *BControl) OutputChannel() <-chan *client.BResponse {
	return c.outputChan
}

// Done should be used by a goroutine to indicate it finished processing
func (c *BControl) Done() {
	defer c.wg.Done()
	c.semaphore.release()
}

// RunWhenAvailable blocks waiting for a execution slot
func (c *BControl) RunWhenAvailable() {
	c.semaphore.acquire()
}
