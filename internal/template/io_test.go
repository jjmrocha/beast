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

package template

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	templatedEndpoint = "http://someendpoint.pt/{{ .RequestID }}"
	templatedBody     = "{\"id\": {{ .RequestID }}, \"value\": \"{{ .Data.A }}\"}"
	staticBody        = "{\"id\": 1, \"value\": \"any\"}"
	userAgent         = "Beast/1"
	headerKey1        = "Key_1"
	headerValue1      = "Value_1"
	userAgentHeader   = "User-Agent"
	basicEndpoint     = "http://www.google.pt"
)

var jsonContentType = []Header{{Key: "Content-Type", Value: "application/json"}}

func TestRead(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Template
	}{
		{
			name:     "basic GET in JSON",
			input:    "../../testdata/basic_get.json",
			expected: &Template{Method: http.MethodGet, Endpoint: basicEndpoint},
		},
		{
			name:     "basic GET in YAML",
			input:    "../../testdata/basic_get.yaml",
			expected: &Template{Method: http.MethodGet, Endpoint: basicEndpoint},
		},
		{
			name:     "basic POST in JSON",
			input:    "../../testdata/basic_post.json",
			expected: &Template{Method: http.MethodPost, Endpoint: "http://someendpoint.pt", Headers: jsonContentType, Body: staticBody},
		},
		{
			name:     "basic POST in YAML",
			input:    "../../testdata/basic_post.yaml",
			expected: &Template{Method: http.MethodPost, Endpoint: "http://someendpoint.pt", Headers: jsonContentType, Body: staticBody},
		},
		{
			name:     "templated POST in JSON",
			input:    "../../testdata/template_post.json",
			expected: &Template{Method: http.MethodPost, Endpoint: templatedEndpoint, Headers: jsonContentType, Body: templatedBody},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			input := tc.input
			// when
			result, err := Read(input)
			// then
			require.NoError(t, err)
			assert.Equal(t, tc.expected.Method, result.Method)
			assert.Equal(t, tc.expected.Endpoint, result.Endpoint)
			assert.ElementsMatch(t, tc.expected.Headers, result.Headers)
			assert.Equal(t, tc.expected.Body, result.Body)
		})
	}

	t.Run("JSON with an external body file", func(t *testing.T) {
		// given
		t.Chdir("..")
		input := "../testdata/template_post_external.json"
		// when
		result, err := Read(input)
		// then
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, result.Method)
		assert.Equal(t, templatedEndpoint, result.Endpoint)
		assert.ElementsMatch(t, jsonContentType, result.Headers)
		assert.Equal(t, templatedBody, result.Body)
	})

	t.Run("a missing file returns an error", func(t *testing.T) {
		// given
		input := filepath.Join(t.TempDir(), "missing.yaml")
		// when
		result, err := Read(input)
		// then
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, result)
	})

	t.Run("a missing @file body returns an error", func(t *testing.T) {
		// given
		dir := t.TempDir()
		input := filepath.Join(dir, "request.json")
		content := `{"method": "POST", "url": "http://a.pt", "body": "@` + filepath.Join(dir, "missing.json") + `"}`
		require.NoError(t, os.WriteFile(input, []byte(content), 0o600))
		// when
		result, err := Read(input)
		// then
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, result)
	})

	t.Run("templated POST in YAML", func(t *testing.T) {
		// given
		input := "../../testdata/template_post.yaml"
		// when
		result, err := Read(input)
		// then
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, result.Method)
		assert.Equal(t, templatedEndpoint, result.Endpoint)
		assert.ElementsMatch(t, jsonContentType, result.Headers)
		assert.Contains(t, result.Body, "\"id\": {{ .RequestID }}")
		assert.Contains(t, result.Body, "\"value\": \"{{ .Data.A }}\"")
	})
}

func TestWrite(t *testing.T) {
	tmpl := &Template{
		Method:   "PUT",
		Endpoint: templatedEndpoint,
		Headers:  []Header{{Key: userAgentHeader, Value: userAgent}, {Key: headerKey1, Value: headerValue1}},
		Body:     templatedBody,
	}

	t.Run("JSON file reads back as the same template", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "request.json")
		// when
		Write(fileName, tmpl)
		// then
		result := mustRead(t, fileName)
		assert.Equal(t, tmpl, result)
	})

	t.Run("JSON file uses the v2.6.1 key names", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "request.json")
		// when
		Write(fileName, tmpl)
		// then
		result := mustReadFile(t, fileName)
		for _, key := range []string{`"method"`, `"url"`, `"headers"`, `"key"`, `"value"`, `"body"`} {
			assert.Contains(t, result, key)
		}
	})

	for _, name := range []string{"request.json", "request.yaml"} {
		t.Run("creates "+name+" readable only by its owner", func(t *testing.T) {
			// given
			fileName := filepath.Join(t.TempDir(), name)
			// when
			Write(fileName, tmpl)
			// then
			info, err := os.Stat(fileName)
			require.NoError(t, err)
			expected := os.FileMode(0o600)
			assert.Equal(t, expected, info.Mode().Perm())
		})
	}

	t.Run("YAML file reads back as the same template", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "request.yaml")
		// when
		Write(fileName, tmpl)
		// then
		result := mustRead(t, fileName)
		assert.Equal(t, tmpl.Method, result.Method)
		assert.Equal(t, tmpl.Endpoint, result.Endpoint)
		assert.ElementsMatch(t, tmpl.Headers, result.Headers)
		assert.Equal(t, tmpl.Body, result.Body)
	})

	t.Run("YAML file uses the v2.6.1 key names", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "request.yaml")
		// when
		Write(fileName, tmpl)
		// then
		result := mustReadFile(t, fileName)
		for _, key := range []string{"method:", "endpoint:", "headers:", "User-Agent:", "request-body:"} {
			assert.Contains(t, result, key)
		}
	})
}

func TestIsJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "upper case extension", input: "SOME_FILE.JSON", expected: true},
		{name: "json only in a folder name", input: "/json/not_ajson", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			input := tc.input
			// when
			result := isJSON(input)
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func mustReadFile(t testing.TB, fileName string) string {
	t.Helper()
	data, err := os.ReadFile(fileName)
	require.NoError(t, err)
	return string(data)
}

func mustRead(t testing.TB, fileName string) *Template {
	t.Helper()
	tmpl, err := Read(fileName)
	if err != nil {
		t.Fatalf("Read(%s): unexpected error: %v", fileName, err)
	}
	return tmpl
}
