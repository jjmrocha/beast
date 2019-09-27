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
package beast

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func HttpClient() *http.Client {
	transport := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives:  false,
		MaxConnsPerHost:    0,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}
}

func Convert(request *BRequest) (*http.Request, error) {
	req, err := http.NewRequest(request.Method, request.Endpoint, bodyReader(request.Body))
	if err != nil {
		return nil, err
	}

	for _, header := range request.Headers {
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

func Execute(client *http.Client, req *http.Request) *BResponse {
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Error executing request '%s %s': %v\n", req.Method, req.URL.String(), err)
		return &BResponse{
			StatusCode: -1,
			Duration:   duration,
		}
	}

	return &BResponse{
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}
