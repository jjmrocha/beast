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

// Package control provides functions that perform the paralel execution of requests
package control

import (
	"log"
	"sync"
	"time"

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/data"
	"github.com/jjmrocha/beast/template"
)

// Control is used to control the execution of multiple goroutines
type Control struct {
	wg                 sync.WaitGroup
	generatorChannel   chan *template.Generator
	outputChannel      chan *client.Response
	requestCount       int
	executionDuration  int
	concurrentRoutines int
}

// New creates a control.Controle
func New(nRequests, tDuration, nParallel int) *Control {
	ctrl := &Control{
		generatorChannel:   make(chan *template.Generator, nParallel),
		outputChannel:      make(chan *client.Response, nParallel),
		requestCount:       nRequests,
		executionDuration:  tDuration,
		concurrentRoutines: nParallel,
	}
	ctrl.wg.Add(nParallel)

	go func() {
		ctrl.wg.Wait()
		close(ctrl.outputChannel)
	}()

	return ctrl
}

// OutputChannel returns a channel for receiving the client.Response sent using Push
func (c *Control) OutputChannel() <-chan *client.Response {
	return c.outputChannel
}

// AsyncExecute creates the goroutines and start the test execution
func (c *Control) AsyncExecute(httpClient *client.Client, tmplc *template.CompiledTemplate, rows *data.Data) {
	go c.createGenerators(tmplc, rows)

	for i := 0; i < c.concurrentRoutines; i++ {
		requestChannel := make(chan *client.Request)
		go c.makeRequest(requestChannel)
		go c.executeRequest(requestChannel, httpClient)
	}
}

// Consts
var emptyRecord = data.NewRecord()
var generateError = client.Response{
	Timestamp:  time.Now(),
	StatusCode: -100,
}

func (c *Control) createGenerators(tmplc *template.CompiledTemplate, rows *data.Data) {
	requestCount := c.requestCount
	if requestCount == 0 {
		requestCount = int(^uint(0) >> 1) // Max int value
	}

	duration := c.executionDuration
	if duration == 0 {
		duration = 31536000 // One year
	}

	nextRecord := func() *data.Record {
		if rows == nil {
			return &emptyRecord
		}

		return rows.Next()
	}

	defer close(c.generatorChannel)
	timeout := time.After(time.Duration(duration) * time.Second)

	for i := 1; i <= requestCount; i++ {
		generator := &template.Generator{
			Data:     nextRecord(),
			RecordID: i,
			Template: tmplc,
		}
		select {
		case c.generatorChannel <- generator:
		case <-timeout:
			return
		}
	}
}

func (c *Control) makeRequest(requestChannel chan<- *client.Request) {
	defer close(requestChannel)

	for generator := range c.generatorChannel {
		req, err := generator.Request()
		if err != nil {
			log.Printf("Error generating request for %s: %v\n", generator.Log(), err)
			c.outputChannel <- &generateError
			continue
		}

		requestChannel <- req
	}
}

func (c *Control) executeRequest(requestChannel <-chan *client.Request, httpClient *client.Client) {
	defer c.wg.Done()

	for req := range requestChannel {
		c.outputChannel <- httpClient.Execute(req)
	}
}
