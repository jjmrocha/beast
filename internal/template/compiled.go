/*
 * Copyright 2020 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
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
	"bytes"
	txt "text/template"

	"github.com/jjmrocha/beast/internal/data"
)

type compiledHeader struct {
	key   string
	value *txt.Template
}

// CompiledTemplate is the compiled version of a Template
type CompiledTemplate struct {
	method   string
	endpoint *txt.Template
	headers  []compiledHeader
	body     *txt.Template
}

type requestValues struct {
	RequestID int
	Data      *data.Record
}

func (c *CompiledTemplate) executeTemplate(requestID int, record *data.Record) (*Template, error) {
	values := requestValues{RequestID: requestID, Data: record}

	endpoint, err := render(c.endpoint, values)
	if err != nil {
		return nil, err
	}

	headers, err := c.renderHeaders(values)
	if err != nil {
		return nil, err
	}

	body, err := render(c.body, values)
	if err != nil {
		return nil, err
	}

	return &Template{Method: c.method, Endpoint: endpoint, Headers: headers, Body: body}, nil
}

func (c *CompiledTemplate) renderHeaders(values requestValues) ([]Header, error) {
	headers := make([]Header, 0, len(c.headers))
	for _, header := range c.headers {
		value, err := render(header.value, values)
		if err != nil {
			return nil, err
		}
		headers = append(headers, Header{Key: header.key, Value: value})
	}
	return headers, nil
}

func render(tmpl *txt.Template, values requestValues) (string, error) {
	if tmpl == nil {
		return "", nil
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, values); err != nil {
		return "", err
	}

	return rendered.String(), nil
}
