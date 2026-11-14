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

// Package template provide functions to read, write templates and generate requests
package template

import (
	"io"
	"net/http"
	"strings"
	"text/template"
)

// Header represents an HTTP Header template
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Template represents an HTTP request template
type Template struct {
	Method   string   `json:"method"`
	Endpoint string   `json:"url"`
	Headers  []Header `json:"headers,omitempty"`
	Body     string   `json:"body,omitempty"`
}

// Compile returns a compiled version of the template
func (t *Template) Compile() (*CompiledTemplate, error) {
	endpoint, err := parse("endpoint", t.Endpoint)
	if err != nil {
		return nil, err
	}

	headers, err := compileHeaders(t.Headers)
	if err != nil {
		return nil, err
	}

	body, err := parse("body", t.Body)
	if err != nil {
		return nil, err
	}

	return &CompiledTemplate{method: t.Method, endpoint: endpoint, headers: headers, body: body}, nil
}

func compileHeaders(headers []Header) ([]compiledHeader, error) {
	compiled := make([]compiledHeader, 0, len(headers))
	for _, header := range headers {
		value, err := parse("headerValue", header.Value)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, compiledHeader{key: header.Key, value: value})
	}
	return compiled, nil
}

func parse(name, text string) (*template.Template, error) {
	if text == "" {
		return nil, nil
	}

	return template.New(name).Parse(text)
}

func (t *Template) request() (*http.Request, error) {
	req, err := http.NewRequest(t.Method, t.Endpoint, bodyReader(t.Body))
	if err != nil {
		return nil, err
	}

	for _, header := range t.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	return req, nil
}

func bodyReader(body string) io.Reader {
	if body == "" {
		return nil
	}

	return strings.NewReader(body)
}
