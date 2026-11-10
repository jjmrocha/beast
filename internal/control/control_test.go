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
	"testing"
	"time"

	"github.com/jjmrocha/beast/internal/client"
	"github.com/jjmrocha/beast/internal/config"
	"github.com/jjmrocha/beast/internal/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsyncExecute(t *testing.T) {
	t.Run("records each render failure as its own result with its own time", func(t *testing.T) {
		// given
		tmpl := &template.Template{Method: "GET", Endpoint: "http://someendpoint.pt/{{ index .Data 1 }}"}
		tmplc, err := tmpl.Compile()
		require.NoError(t, err)
		ctrl := New(Limits{Requests: 2, Workers: 1})
		start := time.Now()
		// when
		ctrl.AsyncExecute(client.NewClient(config.Default(), 1), tmplc, nil)
		var results []*client.Response
		for response := range ctrl.OutputChannel() {
			results = append(results, response)
		}
		// then
		require.Len(t, results, 2)
		assert.NotSame(t, results[0], results[1])
		for _, result := range results {
			assert.Equal(t, client.StatusGenerationError, result.StatusCode)
			assert.False(t, result.Timestamp.Before(start), "timestamp %v is before the run started at %v", result.Timestamp, start)
		}
	})
}
