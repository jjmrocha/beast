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

package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jjmrocha/beast/config"
	"github.com/stretchr/testify/assert"
)

const (
	clientError       = -500
	timeoutError      = -400
	redirectPath      = "/old"
	redirectTarget    = "/new"
	slowHandlerMaxRun = 3 * time.Second
)

func TestNewClient(t *testing.T) {
	t.Run("times out when the server is slower than request-timeout", func(t *testing.T) {
		// given
		server := newServer(t, func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-r.Context().Done():
			case <-time.After(slowHandlerMaxRun):
			}
		})
		cfg := config.Default()
		cfg.RequestTimeout = 1
		c := NewClient(cfg, 1)
		request := mustNewRequest(t, server.URL)
		// when
		result := c.Execute(request)
		// then
		assert.Equal(t, timeoutError, result.StatusCode)
		assert.GreaterOrEqual(t, result.Duration, time.Second)
	})

	t.Run("does not follow redirects by default", func(t *testing.T) {
		// given
		server := newRedirectServer(t)
		c := NewClient(config.Default(), 1)
		request := mustNewRequest(t, server.URL+redirectPath)
		// when
		result := c.Execute(request)
		// then
		assert.Equal(t, http.StatusFound, result.StatusCode)
	})

	t.Run("follows redirects when disable-redirects is false", func(t *testing.T) {
		// given
		server := newRedirectServer(t)
		cfg := config.Default()
		cfg.DisableRedirects = false
		c := NewClient(cfg, 1)
		request := mustNewRequest(t, server.URL+redirectPath)
		// when
		result := c.Execute(request)
		// then
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})

	t.Run("rejects a self-signed certificate by default", func(t *testing.T) {
		// given
		server := newTLSServer(t)
		c := NewClient(config.Default(), 1)
		request := mustNewRequest(t, server.URL)
		// when
		result := c.Execute(request)
		// then
		assert.Equal(t, clientError, result.StatusCode)
	})

	t.Run("accepts a self-signed certificate when disable-certificate-check is true", func(t *testing.T) {
		// given
		server := newTLSServer(t)
		cfg := config.Default()
		cfg.DisableCertificateCheck = true
		c := NewClient(cfg, 1)
		request := mustNewRequest(t, server.URL)
		// when
		result := c.Execute(request)
		// then
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})
}

func TestExecute(t *testing.T) {
	t.Run("returns the status code, the request and the duration", func(t *testing.T) {
		// given
		server := newStatusServer(t, http.StatusOK)
		c := NewClient(config.Default(), 1)
		request := mustNewRequest(t, server.URL+"/apps")
		// when
		result := c.Execute(request)
		// then
		assert.Equal(t, http.StatusOK, result.StatusCode)
		expectedRequest := "GET " + server.URL + "/apps"
		assert.Equal(t, expectedRequest, result.Request)
		assert.Positive(t, result.Duration)
	})

	t.Run("returns a non-success status as it is", func(t *testing.T) {
		// given
		server := newStatusServer(t, http.StatusServiceUnavailable)
		c := NewClient(config.Default(), 1)
		request := mustNewRequest(t, server.URL)
		// when
		result := c.Execute(request)
		// then
		assert.Equal(t, http.StatusServiceUnavailable, result.StatusCode)
	})

	t.Run("returns -500 when the server cannot be reached", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.NotFoundHandler())
		url := server.URL
		server.Close()
		c := NewClient(config.Default(), 1)
		request := mustNewRequest(t, url)
		// when
		result := c.Execute(request)
		// then
		assert.Equal(t, clientError, result.StatusCode)
		assert.Equal(t, "GET "+url, result.Request)
	})
}

func mustNewRequest(t testing.TB, url string) *Request {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("http.NewRequest(%s): unexpected error: %v", url, err)
	}
	return BuildRequest(req)
}

func newServer(t testing.TB, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

func newStatusServer(t testing.TB, status int) *httptest.Server {
	t.Helper()
	return newServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	})
}

func newRedirectServer(t testing.TB) *httptest.Server {
	t.Helper()
	return newServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == redirectPath {
			http.Redirect(w, r, redirectTarget, http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func newTLSServer(t testing.TB) *httptest.Server {
	t.Helper()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	return server
}
