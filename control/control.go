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
package control

import (
	"sync"

	"github.com/jjmrocha/beast/client"
)

type BControl struct {
	wg         sync.WaitGroup
	semaphore  client.Semaphore
	outputChan chan *client.BResponse
}

func New(nRequests, nParallel int) *BControl {
	ctrl := &BControl{
		semaphore:  client.NewSemaphore(nParallel),
		outputChan: make(chan *client.BResponse, nRequests),
	}
	ctrl.wg.Add(nRequests)

	return ctrl
}

func (c *BControl) Send(response *client.BResponse) {
	c.outputChan <- response
	c.semaphore.Release()
}

func (c *BControl) CloseWhenDone() {
	c.wg.Wait()
	close(c.outputChan)
}

func (c *BControl) Output() <-chan *client.BResponse {
	return c.outputChan
}

func (c *BControl) Done() {
	c.wg.Done()
}

func (c *BControl) WaitToStart() {
	c.semaphore.Acquire()
}
