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
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jjmrocha/beast/config"
)

// This interface exists to allow the mocking of the client
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// BClient represents an HTTP client
type BClient struct {
	native httpClient
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
		native: client,
	}
}

// Execute executes the request measuring the time taken to execute and return a BResponse
func (c *BClient) Execute(request *BRequest) *BResponse {
	start := time.Now()
	resp, err := c.native.Do(request.native)
	duration := time.Since(start)

	if err != nil {
		var urlErr *url.Error
		if ok := errors.As(err, &urlErr); ok && urlErr.Timeout() {
			return &BResponse{
				StatusCode: -400,
				Duration:   duration,
			}
		}

		log.Printf("Error executing request '%v': %v\n", request, err)
		return &BResponse{
			StatusCode: -500,
			Duration:   duration,
		}
	}

	defer resp.Body.Close()

	return &BResponse{
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}
