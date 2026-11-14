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

package command

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const runTimeout = 10 * time.Second

func TestRun(t *testing.T) {
	t.Run("prints the test settings, then the stats", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		t.Cleanup(server.Close)
		templateFile := writeTemplate(t, server.URL)
		opts := RunOptions{Requests: 3, Concurrency: 2, Template: templateFile}
		var out strings.Builder
		// when
		err := Run(t.Context(), &out, opts)
		// then
		assert.NoError(t, err)
		result := out.String()
		expectedSettings := "===== Test =====\n" +
			"Request template: " + templateFile + "\n" +
			"Number of requests: 3\n" +
			"Number of concurrent requests: 2\n" +
			"===== Preparing =====\n" +
			"- Loading request template\n" +
			"===== Executing =====\n" +
			"===== Stats =====\n" +
			"Executed requests: 3\n"
		assert.True(t, strings.HasPrefix(result, "===== System =====\n"), "run output: %q", result)
		assert.Contains(t, result, expectedSettings)
	})

	t.Run("returns an error when a file cannot be used", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "missing")
		validTemplate := writeTemplate(t, "http://127.0.0.1:1/")
		tests := []struct {
			name            string
			input           RunOptions
			expectedCause   error
			expectedMessage string
		}{
			{
				name:            "a missing config file",
				input:           RunOptions{Requests: 1, Concurrency: 1, Template: validTemplate, Config: missing},
				expectedCause:   os.ErrNotExist,
				expectedMessage: "reading config file " + missing + ": ",
			},
			{
				name:            "a missing template file",
				input:           RunOptions{Requests: 1, Concurrency: 1, Template: missing},
				expectedCause:   os.ErrNotExist,
				expectedMessage: "reading template file " + missing + ": ",
			},
			{
				name:            "a template that does not compile",
				input:           RunOptions{Requests: 1, Concurrency: 1, Template: writeTemplate(t, "http://a.pt/{{")},
				expectedMessage: "compiling template: ",
			},
			{
				name:            "a missing data file",
				input:           RunOptions{Requests: 1, Concurrency: 1, Template: validTemplate, Data: missing},
				expectedCause:   os.ErrNotExist,
				expectedMessage: "reading data file " + missing + ": ",
			},
			{
				name:            "an output file in a missing folder",
				input:           RunOptions{Requests: 1, Concurrency: 1, Template: validTemplate, Output: filepath.Join(missing, "output.csv")},
				expectedCause:   os.ErrNotExist,
				expectedMessage: "creating output file " + filepath.Join(missing, "output.csv") + ": ",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given
				opts := tc.input
				// when
				err := Run(t.Context(), io.Discard, opts)
				// then
				assert.ErrorContains(t, err, tc.expectedMessage)
				if tc.expectedCause != nil {
					assert.ErrorIs(t, err, tc.expectedCause)
				}
			})
		}
	})

	t.Run("sends no request when the output file cannot be created", func(t *testing.T) {
		// given
		var received atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			received.Add(1)
		}))
		t.Cleanup(server.Close)
		opts := RunOptions{Requests: 5, Concurrency: 1, Template: writeTemplate(t, server.URL), Output: filepath.Join(t.TempDir(), "missing", "output.csv")}
		// when
		err := Run(t.Context(), io.Discard, opts)
		// then
		assert.Error(t, err)
		assert.Zero(t, received.Load())
	})

	t.Run("writes a complete output file when the run is cancelled", func(t *testing.T) {
		// given
		arrived := make(chan struct{}, 16)
		released := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case arrived <- struct{}{}:
			default:
			}
			select {
			case <-released:
			case <-r.Context().Done():
			}
		}))
		t.Cleanup(server.Close)
		templateFile := writeTemplate(t, server.URL)
		opts := RunOptions{Duration: 60, Concurrency: 2, Template: templateFile, Output: filepath.Join(t.TempDir(), "output.csv")}
		ctx, cancel := context.WithCancel(t.Context())
		finished := make(chan struct{})
		go func() {
			defer close(finished)
			assert.NoError(t, Run(ctx, io.Discard, opts))
		}()
		<-arrived
		<-arrived
		// when
		cancel()
		close(released)
		// then
		select {
		case <-finished:
		case <-time.After(runTimeout):
			t.Fatalf("Run did not return %v after the cancel", runTimeout)
		}
		file, err := os.Open(opts.Output)
		require.NoError(t, err)
		defer func() { _ = file.Close() }()
		lines, err := csv.NewReader(file).ReadAll()
		require.NoError(t, err)
		expectedLines := 3
		assert.Len(t, lines, expectedLines, "the header and one line for each request in flight")
	})
}

func writeTemplate(t testing.TB, url string) string {
	t.Helper()
	fileName := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(fileName, []byte(`{"method": "GET", "url": "`+url+`"}`), 0o600); err != nil {
		t.Fatalf("writing template %s: %v", fileName, err)
	}
	return fileName
}
