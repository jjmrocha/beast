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

package request

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/data"
)

type cHeader struct {
	key   string
	value *template.Template
}

type cRequest struct {
	method   string
	endpoint *template.Template
	headers  []cHeader
	body     *template.Template
}

func (c *cRequest) executeTemplate(requestID int, record *data.Record) (*TRequest, error) {
	var context = struct {
		RequestID int
		Data      *data.Record
	}{
		RequestID: requestID,
		Data:      record,
	}

	tRequest := TRequest{
		Method:  c.method,
		Headers: make([]THeader, 0, len(c.headers)),
	}

	var endpoint bytes.Buffer
	if err := c.endpoint.Execute(&endpoint, context); err != nil {
		return nil, err
	}
	tRequest.Endpoint = endpoint.String()

	for _, header := range c.headers {
		var headerValue bytes.Buffer
		if err := header.value.Execute(&headerValue, context); err != nil {
			return nil, err
		}
		tHeader := THeader{
			Key:   header.key,
			Value: headerValue.String(),
		}
		tRequest.Headers = append(tRequest.Headers, tHeader)
	}

	if c.body != nil {
		var body bytes.Buffer
		if err := c.body.Execute(&body, context); err != nil {
			return nil, err
		}
		tRequest.Body = body.String()
	}

	return &tRequest, nil
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
