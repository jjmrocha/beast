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

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/data"
)

// THeader represents an HTTP Header template
type THeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TRequest represents an HTTP request template
type TRequest struct {
	Method   string    `json:"method"`
	Endpoint string    `json:"url"`
	Headers  []THeader `json:"headers,omitempty"`
	Body     string    `json:"body,omitempty"`
}

func (t *TRequest) request() (*client.BRequest, error) {
	req, err := http.NewRequest(t.Method, t.Endpoint, bodyReader(t.Body))
	if err != nil {
		return nil, err
	}

	for _, header := range t.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	return client.MakeRequest(req), nil
}

func bodyReader(body string) io.Reader {
	if body == "" {
		return nil
	}

	return strings.NewReader(body)
}

func (t *TRequest) compile() (*cRequest, error) {
	tEndpoint, err := template.New("endpoint").Parse(t.Endpoint)
	if err != nil {
		return nil, err
	}

	tHeaders := make([]cHeader, 0, len(t.Headers))
	for _, header := range t.Headers {
		tValue, err := template.New("headerValue").Parse(header.Value)
		if err != nil {
			return nil, err
		}
		cHeader := cHeader{
			key:   header.Key,
			value: tValue,
		}
		tHeaders = append(tHeaders, cHeader)
	}

	var tBody *template.Template
	if t.Body != "" {
		tBody, err = template.New("body").Parse(t.Body)
		if err != nil {
			return nil, err
		}
	}

	cRequest := cRequest{
		method:   t.Method,
		endpoint: tEndpoint,
		headers:  tHeaders,
		body:     tBody,
	}
	return &cRequest, nil
}

// BuildGenerators creates a slice with requests generators
func (t *TRequest) BuildGenerators(nRequests int, data *data.Data) ([]*Generator, error) {
	cRequest, err := t.compile()
	if err != nil {
		return nil, err
	}

	generators := make([]*Generator, 0, nRequests)

	for i := 1; i <= nRequests; i++ {
		generator := &Generator{
			data:     nextRecord(data),
			recordID: i,
			template: cRequest,
		}
		generators = append(generators, generator)
	}

	return generators, nil
}

func nextRecord(data *data.Data) *data.Record {
	if data == nil {
		return &emptyRecord
	}

	return data.Next()
}
