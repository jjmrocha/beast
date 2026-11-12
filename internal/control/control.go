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

// Package control provides functions that perform the parallel execution of requests
package control

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jjmrocha/beast/internal/client"
	"github.com/jjmrocha/beast/internal/data"
	"github.com/jjmrocha/beast/internal/template"
)

type Limits struct {
	Requests int
	Duration int
	Workers  int
}

func Start(ctx context.Context, limits Limits, httpClient *client.Client, tmplc *template.CompiledTemplate, rows *data.Data) <-chan *client.Response {
	ctx, cancel := withDuration(ctx, limits.Duration)

	generators := make(chan *template.Generator, limits.Workers)
	results := make(chan *client.Response, limits.Workers)
	go handOut(ctx, limits.Requests, tmplc, rows, generators)

	var wg sync.WaitGroup
	for range limits.Workers {
		requests := make(chan *http.Request)
		wg.Go(func() { render(ctx, generators, requests, results) })
		wg.Go(func() { send(ctx, requests, httpClient, results) })
	}

	go func() {
		wg.Wait()
		cancel()
		close(results)
	}()

	return results
}

func withDuration(ctx context.Context, seconds int) (context.Context, context.CancelFunc) {
	if seconds == 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, time.Duration(seconds)*time.Second)
}

var emptyRecord = data.NewRecord()

func handOut(ctx context.Context, requests int, tmplc *template.CompiledTemplate, rows *data.Data, generators chan<- *template.Generator) {
	defer close(generators)

	nextRecord := func() *data.Record {
		if rows == nil {
			return &emptyRecord
		}

		return rows.Next()
	}

	for i := 1; requests == 0 || i <= requests; i++ {
		generator := &template.Generator{
			Data:     nextRecord(),
			RecordID: i,
			Template: tmplc,
		}
		select {
		case generators <- generator:
		case <-ctx.Done():
			return
		}
	}
}

func render(ctx context.Context, generators <-chan *template.Generator, requests chan<- *http.Request, results chan<- *client.Response) {
	defer close(requests)

	for generator := range generators {
		if ctx.Err() != nil {
			return
		}

		req, err := generator.Request()
		if err != nil {
			log.Printf("Error generating request for %s: %v\n", generator.Log(), err)
			results <- &client.Response{
				Timestamp:  time.Now(),
				StatusCode: client.StatusGenerationError,
			}
			continue
		}

		select {
		case requests <- req:
		case <-ctx.Done():
			return
		}
	}
}

func send(ctx context.Context, requests <-chan *http.Request, httpClient *client.Client, results chan<- *client.Response) {
	for req := range requests {
		if ctx.Err() != nil {
			continue
		}

		results <- httpClient.Execute(req)
	}
}
