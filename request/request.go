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
	"fmt"

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/data"
)

// Generator generates BRequests
type Generator struct {
	final    *client.BRequest
	template *cRequest
	recordID int
	data     *data.Record
}

// Request uses that template and a record and returns a BRequests
func (g *Generator) Request() (*client.BRequest, error) {
	if g.final != nil {
		return g.final, nil
	}

	dRequest, err := g.template.executeTemplate(g.recordID, g.data)
	if err != nil {
		return nil, err
	}

	bRequest, err := dRequest.request()
	if err != nil {
		return nil, err
	}

	return bRequest, nil
}

// Log generates a log message for the request
func (g *Generator) Log() string {
	return fmt.Sprintf("requestId: %v and data: %v", g.recordID, g.data)
}

// CreateRequests creates a slice with requests generators
func (t *TRequest) CreateRequests(nRequests int, data *data.Data) ([]*Generator, error) {
	if data == nil {
		return staticRequests(nRequests, t)
	}

	return dynamicRequests(nRequests, t, data)
}

func staticRequests(nRequests int, tRequest *TRequest) ([]*Generator, error) {
	bRequest, err := tRequest.request()
	if err != nil {
		return nil, err
	}

	generator := &Generator{final: bRequest}
	generators := make([]*Generator, 0, nRequests)

	for i := 0; i < nRequests; i++ {
		generators = append(generators, generator)
	}

	return generators, nil
}

func dynamicRequests(nRequests int, tRequest *TRequest, data *data.Data) ([]*Generator, error) {
	cRequest, err := tRequest.compile()
	if err != nil {
		return nil, err
	}

	generators := make([]*Generator, 0, nRequests)

	for i := 1; i <= nRequests; i++ {
		generator := &Generator{
			data:     data.Next(),
			recordID: i,
			template: cRequest,
		}
		generators = append(generators, generator)
	}

	return generators, nil
}
