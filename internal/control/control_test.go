/*
 * Copyright 2026 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
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
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jjmrocha/beast/internal/client"
	"github.com/jjmrocha/beast/internal/config"
	"github.com/jjmrocha/beast/internal/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const resultsTimeout = 10 * time.Second

func TestStart(t *testing.T) {
	t.Run("sends the requested number of requests and closes the channel", func(t *testing.T) {
		// given
		server := newHoldingServer(t)
		server.release()
		limits := Limits{Requests: 5, Workers: 2}
		// when
		results := Start(t.Context(), limits, client.NewClient(config.Default(), 2), mustCompile(t, server.URL), nil)
		// then
		responses := collect(t, results)
		require.Len(t, responses, 5)
		for _, response := range responses {
			assert.Equal(t, http.StatusOK, response.StatusCode)
		}
		assert.Equal(t, int64(5), server.hits.Load())
	})

	t.Run("cancelling the context lets requests in flight finish and sends no new request", func(t *testing.T) {
		// given
		server := newHoldingServer(t)
		ctx, cancel := context.WithCancel(t.Context())
		results := Start(ctx, Limits{Workers: 2}, client.NewClient(config.Default(), 2), mustCompile(t, server.URL), nil)
		server.waitForArrivals(t, 2)
		// when
		cancel()
		server.release()
		// then
		responses := collect(t, results)
		require.Len(t, responses, 2)
		for _, response := range responses {
			assert.Equal(t, http.StatusOK, response.StatusCode)
		}
		assert.Equal(t, int64(2), server.hits.Load())
	})

	t.Run("drops a rendered request that is still waiting when the deadline passes", func(t *testing.T) {
		// given
		server := newHoldingServer(t)
		time.AfterFunc(1500*time.Millisecond, server.release)
		limits := Limits{Duration: 1, Workers: 1}
		// when
		results := Start(t.Context(), limits, client.NewClient(config.Default(), 1), mustCompile(t, server.URL), nil)
		// then
		responses := collect(t, results)
		require.Len(t, responses, 1)
		assert.Equal(t, http.StatusOK, responses[0].StatusCode)
		assert.Equal(t, int64(1), server.hits.Load())
	})

	t.Run("sends nothing when the context is cancelled before the start", func(t *testing.T) {
		// given
		server := newHoldingServer(t)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		// when
		results := Start(ctx, Limits{Workers: 2}, client.NewClient(config.Default(), 2), mustCompile(t, server.URL), nil)
		// then
		assert.Empty(t, collect(t, results))
		assert.Equal(t, int64(0), server.hits.Load())
	})

	t.Run("records each render failure as its own result with its own time", func(t *testing.T) {
		// given
		tmpl := &template.Template{Method: "GET", Endpoint: "http://someendpoint.pt/{{ index .Data 1 }}"}
		tmplc, err := tmpl.Compile()
		require.NoError(t, err)
		start := time.Now()
		// when
		results := Start(t.Context(), Limits{Requests: 2, Workers: 1}, client.NewClient(config.Default(), 1), tmplc, nil)
		// then
		responses := collect(t, results)
		require.Len(t, responses, 2)
		assert.NotSame(t, responses[0], responses[1])
		for _, response := range responses {
			assert.Equal(t, client.StatusGenerationError, response.StatusCode)
			assert.False(t, response.Timestamp.Before(start), "timestamp %v is before the run started at %v", response.Timestamp, start)
		}
	})
}

type holdingServer struct {
	*httptest.Server
	hits      atomic.Int64
	arrived   chan struct{}
	released  chan struct{}
	closeOnce sync.Once
}

func newHoldingServer(t testing.TB) *holdingServer {
	t.Helper()
	s := &holdingServer{arrived: make(chan struct{}, 16), released: make(chan struct{})}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.hits.Add(1)
		select {
		case s.arrived <- struct{}{}:
		default:
		}
		select {
		case <-s.released:
		case <-r.Context().Done():
		}
	}))
	t.Cleanup(func() {
		s.release()
		s.Close()
	})
	return s
}

func (s *holdingServer) release() {
	s.closeOnce.Do(func() { close(s.released) })
}

func (s *holdingServer) waitForArrivals(t testing.TB, n int) {
	t.Helper()
	for range n {
		select {
		case <-s.arrived:
		case <-time.After(resultsTimeout):
			t.Fatalf("only some of the %d requests reached the server", n)
		}
	}
}

func mustCompile(t testing.TB, url string) *template.CompiledTemplate {
	t.Helper()
	tmpl := &template.Template{Method: http.MethodGet, Endpoint: url}
	tmplc, err := tmpl.Compile()
	if err != nil {
		t.Fatalf("compiling the template for %s: %v", url, err)
	}
	return tmplc
}

func collect(t testing.TB, results <-chan *client.Response) []*client.Response {
	t.Helper()
	var responses []*client.Response
	timeout := time.After(resultsTimeout)
	for {
		select {
		case response, open := <-results:
			if !open {
				return responses
			}
			responses = append(responses, response)
		case <-timeout:
			t.Fatalf("the results channel was still open after %v", resultsTimeout)
			return nil
		}
	}
}
