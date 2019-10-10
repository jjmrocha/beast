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

// Package client provides functions and types for the execution of HTTP requests
package client

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jjmrocha/beast/config"
)

// BRequest represents an HTTP request
type BRequest struct {
	Native *http.Request
}

func (r *BRequest) String() string {
	return fmt.Sprintf("%s %s", r.Native.Method, r.Native.URL)
}

// BResponse contains the status code and the duration taken for execution of a request
type BResponse struct {
	StatusCode int
	Duration   time.Duration
}

func (r *BResponse) String() string {
	return fmt.Sprintf("%v - %v", r.StatusCode, r.Duration)
}

// BClient represents an HTTP client
type BClient struct {
	Native *http.Client
}

// HTTP creates a BClient based on the provided configuration
func HTTP(config *config.Config) *BClient {
	tls := &tls.Config{
		InsecureSkipVerify: config.DisableCertificateCheck,
	}
	transport := &http.Transport{
		DisableCompression: config.DisableCompression,
		DisableKeepAlives:  config.DisableKeepAlives,
		MaxConnsPerHost:    config.MaxConnections,
		TLSClientConfig:    tls,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(time.Second.Nanoseconds() * int64(config.RequestTimeout)),
	}

	if config.DisableRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &BClient{
		Native: client,
	}
}

// Execute executes the request measuring the time taken to execute and return a BResponse
func (c *BClient) Execute(request *BRequest) *BResponse {
	start := time.Now()
	resp, err := c.Native.Do(request.Native)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Error executing request '%v': %v\n", request, err)
		return &BResponse{
			StatusCode: -1,
			Duration:   duration,
		}
	}

	defer resp.Body.Close()

	return &BResponse{
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}
