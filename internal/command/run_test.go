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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const runTimeout = 10 * time.Second

func TestRun(t *testing.T) {
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
		dir := t.TempDir()
		templateFile := filepath.Join(dir, "request.json")
		require.NoError(t, os.WriteFile(templateFile, []byte(`{"method": "GET", "url": "`+server.URL+`"}`), 0o600))
		opts := RunOptions{Duration: 60, Concurrency: 2, Template: templateFile, Output: filepath.Join(dir, "output.csv")}
		ctx, cancel := context.WithCancel(t.Context())
		finished := make(chan struct{})
		go func() {
			defer close(finished)
			Run(ctx, opts)
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
