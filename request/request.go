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
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/jjmrocha/beast/data"

	"github.com/jjmrocha/beast/client"
)

type THeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TRequest struct {
	Method   string    `json:"method"`
	Endpoint string    `json:"url"`
	Headers  []THeader `json:"headers,omitempty"`
	Body     string    `json:"body,omitempty"`
}

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

func (c *cRequest) execute(requestID int, record *data.Record) (*TRequest, error) {
	var context = struct {
		RequestID int
		Data      *data.Record
	}{
		RequestID: requestID,
		Data:      record,
	}

	request := TRequest{
		Method:  c.method,
		Headers: make([]THeader, 0, len(c.headers)),
	}

	var endpoint bytes.Buffer
	if err := c.endpoint.Execute(&endpoint, context); err != nil {
		return nil, err
	}
	request.Endpoint = endpoint.String()

	for _, header := range c.headers {
		var headerValue bytes.Buffer
		if err := header.value.Execute(&headerValue, context); err != nil {
			return nil, err
		}
		h := THeader{
			Key:   header.key,
			Value: headerValue.String(),
		}
		request.Headers = append(request.Headers, h)
	}

	if c.body != nil {
		var body bytes.Buffer
		if err := c.body.Execute(&body, context); err != nil {
			return nil, err
		}
		request.Body = body.String()
	}

	return &request, nil
}

func (t *TRequest) Generate(nRequests int, data *data.Data) ([]*client.BRequest, error) {
	if data == nil {
		return monoRequest(nRequests, t)
	} else {
		return multiRequest(nRequests, t, data)
	}
}

func monoRequest(nRequests int, request *TRequest) ([]*client.BRequest, error) {
	bRequest, err := request.generate()
	if err != nil {
		return nil, err
	}

	requests := make([]*client.BRequest, 0, nRequests)
	for i := 0; i < nRequests; i++ {
		requests = append(requests, bRequest)
	}

	return requests, nil
}

func multiRequest(nRequests int, request *TRequest, data *data.Data) ([]*client.BRequest, error) {
	compiled, err := request.toTemplate()
	if err != nil {
		return nil, err
	}

	requests := make([]*client.BRequest, 0, nRequests)
	for i := 1; i <= nRequests; i++ {
		generated, err := compiled.execute(i, data.Next())
		if err != nil {
			return nil, err
		}

		bRequest, err := generated.generate()
		if err != nil {
			return nil, err
		}

		requests = append(requests, bRequest)
	}

	return requests, nil
}

func (t *TRequest) generate() (*client.BRequest, error) {
	req, err := http.NewRequest(t.Method, t.Endpoint, bodyReader(t.Body))
	if err != nil {
		return nil, err
	}

	for _, header := range t.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	return &client.BRequest{Native: req}, nil
}

func bodyReader(body string) io.Reader {
	if body == "" {
		return nil
	}

	return strings.NewReader(body)
}

func (t *TRequest) toTemplate() (*cRequest, error) {
	tEndpoint, err := template.New("endpoint").Parse(t.Endpoint)
	if err != nil {
		return nil, err
	}

	tHeaders := make([]cHeader, 0)
	for _, header := range t.Headers {
		tValue, err := template.New("headerValue").Parse(header.Value)
		if err != nil {
			return nil, err
		}
		ch := cHeader{
			key:   header.Key,
			value: tValue,
		}
		tHeaders = append(tHeaders, ch)
	}

	var tBody *template.Template
	if t.Body != "" {
		tBody, err = template.New("body").Parse(t.Body)
		if err != nil {
			return nil, err
		}
	}

	cr := cRequest{
		method:   t.Method,
		endpoint: tEndpoint,
		headers:  tHeaders,
		body:     tBody,
	}
	return &cr, nil
}

func Read(fileName string) *TRequest {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading template file %s: %v\n", fileName, err)
	}

	var request TRequest
	json.Unmarshal(data, &request)
	return &request
}

func Write(fileName string, request *TRequest) {
	data, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		log.Printf("Error encoding request %v to JSON: %v\n", request, err)
	}

	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		log.Printf("Error writing template to file %s: %v\n", fileName, err)
	}
}
