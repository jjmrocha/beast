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
			input:    "../testdata/basic_get.json",
			expected: &Template{Method: "GET", Endpoint: "http://www.google.pt"},
		},
		{
			name:     "basic GET in YAML",
			input:    "../testdata/basic_get.yaml",
			expected: &Template{Method: "GET", Endpoint: "http://www.google.pt"},
		},
		{
			name:     "basic POST in JSON",
			input:    "../testdata/basic_post.json",
			expected: &Template{Method: "POST", Endpoint: "http://someendpoint.pt", Headers: jsonContentType, Body: staticBody},
		},
		{
			name:     "basic POST in YAML",
			input:    "../testdata/basic_post.yaml",
			expected: &Template{Method: "POST", Endpoint: "http://someendpoint.pt", Headers: jsonContentType, Body: staticBody},
		},
		{
			name:     "templated POST in JSON",
			input:    "../testdata/template_post.json",
			expected: &Template{Method: "POST", Endpoint: templatedEndpoint, Headers: jsonContentType, Body: templatedBody},
		},
		{
			name:     "JSON with an external body file",
			input:    "../testdata/template_post_external.json",
			expected: &Template{Method: "POST", Endpoint: templatedEndpoint, Headers: jsonContentType, Body: templatedBody},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			input := tc.input
			// when
			result := Read(input)
			// then
			assert.Equal(t, tc.expected.Method, result.Method)
			assert.Equal(t, tc.expected.Endpoint, result.Endpoint)
			assert.ElementsMatch(t, tc.expected.Headers, result.Headers)
			assert.Equal(t, tc.expected.Body, result.Body)
		})
	}

	t.Run("templated POST in YAML", func(t *testing.T) {
		// given
		input := "../testdata/template_post.yaml"
		// when
		result := Read(input)
		// then
		assert.Equal(t, "POST", result.Method)
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
		Headers:  []Header{{Key: "User-Agent", Value: "Beast/1"}},
		Body:     templatedBody,
	}

	t.Run("JSON file reads back as the same template", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "request.json")
		// when
		Write(fileName, tmpl)
		// then
		result := Read(fileName)
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

	t.Run("YAML file reads back as the same template", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "request.yaml")
		// when
		Write(fileName, tmpl)
		// then
		result := Read(fileName)
		assert.Equal(t, tmpl, result)
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

func TestExternalBody(t *testing.T) {
	t.Run("reads the file named after @", func(t *testing.T) {
		// given
		body := "@../testdata/body.json"
		// when
		result, found := externalBody(body)
		// then
		assert.True(t, found)
		assert.Equal(t, templatedBody, result)
	})

	t.Run("ignores a body without @", func(t *testing.T) {
		// given
		body := templatedBody
		// when
		result, found := externalBody(body)
		// then
		assert.False(t, found)
		assert.Empty(t, result)
	})
}

func TestIsJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "relative JSON path", input: "../testdata/basic_get.json", expected: true},
		{name: "relative YAML path", input: "../testdata/basic_get.yaml", expected: false},
		{name: "JSON file name", input: "template_post_external.json", expected: true},
		{name: "YAML file name", input: "template_post_external.yaml", expected: false},
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

func TestToHeaderMap(t *testing.T) {
	t.Run("turns the header list into a map", func(t *testing.T) {
		// given
		headers := []Header{{Key: "Key_1", Value: "Value_1"}, {Key: "Key_2", Value: "Value_2"}}
		// when
		result := toHeaderMap(headers)
		// then
		expected := map[string]string{"Key_1": "Value_1", "Key_2": "Value_2"}
		assert.Equal(t, expected, result)
	})

	t.Run("keeps no headers as nil", func(t *testing.T) {
		// when
		result := toHeaderMap(nil)
		// then
		assert.Nil(t, result)
	})
}

func TestFromHeaderMap(t *testing.T) {
	t.Run("turns the header map into a list", func(t *testing.T) {
		// given
		headers := map[string]string{"Key_1": "Value_1", "Key_2": "Value_2"}
		// when
		result := fromHeaderMap(headers)
		// then
		expected := []Header{{Key: "Key_1", Value: "Value_1"}, {Key: "Key_2", Value: "Value_2"}}
		assert.ElementsMatch(t, expected, result)
	})

	t.Run("keeps no headers as nil", func(t *testing.T) {
		// when
		result := fromHeaderMap(nil)
		// then
		assert.Nil(t, result)
	})
}

func TestToYamlTemplate(t *testing.T) {
	t.Run("full template", func(t *testing.T) {
		// given
		tmpl := exampleTemplate()
		// when
		result := toYamlTemplate(tmpl)
		// then
		assert.Equal(t, exampleYamlTemplate(), result)
	})

	t.Run("template without headers or body", func(t *testing.T) {
		// given
		tmpl := &Template{Method: "GET", Endpoint: "http://www.google.pt"}
		// when
		result := toYamlTemplate(tmpl)
		// then
		expected := &templateY{Method: "GET", Endpoint: "http://www.google.pt"}
		assert.Equal(t, expected, result)
	})
}

func TestFromYamlTemplate(t *testing.T) {
	t.Run("full template", func(t *testing.T) {
		// given
		tmply := exampleYamlTemplate()
		// when
		result := fromYamlTemplate(tmply)
		// then
		assert.Equal(t, exampleTemplate(), result)
	})

	t.Run("template without headers or body", func(t *testing.T) {
		// given
		tmply := &templateY{Method: "GET", Endpoint: "http://www.google.pt"}
		// when
		result := fromYamlTemplate(tmply)
		// then
		expected := &Template{Method: "GET", Endpoint: "http://www.google.pt"}
		assert.Equal(t, expected, result)
	})
}

func exampleTemplate() *Template {
	return &Template{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
		Headers:  []Header{{Key: "User-Agent", Value: "Beast/1"}},
		Body:     "Optional, enter body to send with POST or PUT",
	}
}

func exampleYamlTemplate() *templateY {
	return &templateY{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
		Headers:  map[string]string{"User-Agent": "Beast/1"},
		Body:     "Optional, enter body to send with POST or PUT",
	}
}

func mustReadFile(t testing.TB, fileName string) string {
	t.Helper()
	data, err := os.ReadFile(fileName)
	require.NoError(t, err)
	return string(data)
}
