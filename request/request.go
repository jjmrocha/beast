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
	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/data"
)

// Generate creates a slice with all requests to be executed
func (t *TRequest) Generate(nRequests int, data *data.Data) ([]*client.BRequest, error) {
	if data == nil {
		return staticRequests(nRequests, t)
	}

	return dynamicRequests(nRequests, t, data)
}

func staticRequests(nRequests int, tRequest *TRequest) ([]*client.BRequest, error) {
	bRequest, err := tRequest.request()
	if err != nil {
		return nil, err
	}

	requests := make([]*client.BRequest, 0, nRequests)
	for i := 0; i < nRequests; i++ {
		requests = append(requests, bRequest)
	}

	return requests, nil
}

func dynamicRequests(nRequests int, tRequest *TRequest, data *data.Data) ([]*client.BRequest, error) {
	cRequest, err := tRequest.compile()
	if err != nil {
		return nil, err
	}

	requests := make([]*client.BRequest, 0, nRequests)
	for i := 1; i <= nRequests; i++ {
		dRequest, err := cRequest.executeTemplate(i, data.Next())
		if err != nil {
			return nil, err
		}

		bRequest, err := dRequest.request()
		if err != nil {
			return nil, err
		}

		requests = append(requests, bRequest)
	}

	return requests, nil
}
