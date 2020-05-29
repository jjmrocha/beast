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
	"bytes"
	txt "text/template"

	"github.com/jjmrocha/beast/data"
)

type headerC struct {
	key   string
	value *txt.Template
}

type templateC struct {
	method   string
	endpoint *txt.Template
	headers  []headerC
	body     *txt.Template
}

func (c *templateC) executeTemplate(requestID int, record *data.Record) (*Template, error) {
	var context = struct {
		RequestID int
		Data      *data.Record
	}{
		RequestID: requestID,
		Data:      record,
	}

	tmplf := Template{
		Method:  c.method,
		Headers: make([]Header, 0, len(c.headers)),
	}

	var endpoint bytes.Buffer
	if err := c.endpoint.Execute(&endpoint, context); err != nil {
		return nil, err
	}
	tmplf.Endpoint = endpoint.String()

	for _, header := range c.headers {
		var headerValue bytes.Buffer
		if err := header.value.Execute(&headerValue, context); err != nil {
			return nil, err
		}
		hdf := Header{
			Key:   header.key,
			Value: headerValue.String(),
		}
		tmplf.Headers = append(tmplf.Headers, hdf)
	}

	if c.body != nil {
		var body bytes.Buffer
		if err := c.body.Execute(&body, context); err != nil {
			return nil, err
		}
		tmplf.Body = body.String()
	}

	return &tmplf, nil
}
