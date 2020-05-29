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

// Package template provide functions to read, write templates and generate requests
package template

import (
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/data"
)

// emptyRecord contains a empty record to be used when no data is provided
var emptyRecord = data.NewRecord()

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

func (t *Template) request() (*client.Request, error) {
	req, err := http.NewRequest(t.Method, t.Endpoint, bodyReader(t.Body))
	if err != nil {
		return nil, err
	}

	for _, header := range t.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	return client.BuildRequest(req), nil
}

func bodyReader(body string) io.Reader {
	if body == "" {
		return nil
	}

	return strings.NewReader(body)
}

func (t *Template) compile() (*templateC, error) {
	tEndpoint, err := template.New("endpoint").Parse(t.Endpoint)
	if err != nil {
		return nil, err
	}

	tHeaders := make([]headerC, 0, len(t.Headers))
	for _, header := range t.Headers {
		tValue, err := template.New("headerValue").Parse(header.Value)
		if err != nil {
			return nil, err
		}
		hdc := headerC{
			key:   header.Key,
			value: tValue,
		}
		tHeaders = append(tHeaders, hdc)
	}

	var tBody *template.Template
	if t.Body != "" {
		tBody, err = template.New("body").Parse(t.Body)
		if err != nil {
			return nil, err
		}
	}

	tmplc := templateC{
		method:   t.Method,
		endpoint: tEndpoint,
		headers:  tHeaders,
		body:     tBody,
	}
	return &tmplc, nil
}

// BuildGenerators creates a slice with requests generators
func (t *Template) BuildGenerators(nRequests int, data *data.Data) ([]*Generator, error) {
	tmplc, err := t.compile()
	if err != nil {
		return nil, err
	}

	generators := make([]*Generator, 0, nRequests)

	for i := 1; i <= nRequests; i++ {
		gnt := &Generator{
			data:     nextRecord(data),
			recordID: i,
			template: tmplc,
		}
		generators = append(generators, gnt)
	}

	return generators, nil
}

func nextRecord(data *data.Data) *data.Record {
	if data == nil {
		return &emptyRecord
	}

	return data.Next()
}
