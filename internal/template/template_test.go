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

package template

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const unclosedAction = "{{ .RequestID"

func TestCompile(t *testing.T) {
	t.Run("compiles a valid template", func(t *testing.T) {
		// given
		tmpl := &Template{Method: http.MethodPost, Endpoint: templatedEndpoint, Headers: jsonContentType, Body: templatedBody}
		// when
		result, err := tmpl.Compile()
		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	tests := []struct {
		name  string
		input *Template
	}{
		{
			name:  "bad syntax in the URL",
			input: &Template{Method: http.MethodGet, Endpoint: unclosedAction},
		},
		{
			name:  "bad syntax in a header value",
			input: &Template{Method: http.MethodGet, Endpoint: templatedEndpoint, Headers: []Header{{Key: "X-Id", Value: unclosedAction}}},
		},
		{
			name:  "bad syntax in the body",
			input: &Template{Method: http.MethodPost, Endpoint: templatedEndpoint, Body: unclosedAction},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			tmpl := tc.input
			// when
			result, err := tmpl.Compile()
			// then
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}
