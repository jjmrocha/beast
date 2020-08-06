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

// Package client provides functions and types for the execution of HTTP requests
package client

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
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

// Client represents an HTTP client
type Client struct {
	native httpClient
}

// NewClient creates a client.Client based on the provided configuration
func NewClient(cfg *config.Config, parallelConns int) *Client {
	tls := &tls.Config{
		InsecureSkipVerify: cfg.DisableCertificateCheck,
	}
	maxIdleConns := cfg.GetMaxIdleConnections(parallelConns)
	transport := &http.Transport{
		DisableCompression:  cfg.DisableCompression,
		DisableKeepAlives:   cfg.DisableKeepAlives,
		MaxConnsPerHost:     cfg.MaxConnections,
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConns,
		TLSClientConfig:     tls,
	}
	native := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(time.Second.Nanoseconds() * int64(cfg.RequestTimeout)),
	}

	if cfg.DisableRedirects {
		native.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &Client{
		native: native,
	}
}

// Execute executes the request measuring the time taken to execute and return a client.Response
func (c *Client) Execute(request *Request) *Response {
	start := time.Now()
	resp, err := c.native.Do(request.native)
	duration := time.Since(start)

	if err != nil {
		var urlErr *url.Error
		if ok := errors.As(err, &urlErr); ok && urlErr.Timeout() {
			return &Response{
				Timestamp:  start,
				Request:    request.String(),
				StatusCode: -400,
				Duration:   duration,
			}
		}

		log.Printf("Error executing request '%v': %v\n", request, err)
		return &Response{
			Timestamp:  start,
			Request:    request.String(),
			StatusCode: -500,
			Duration:   duration,
		}
	}

	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	return &Response{
		Timestamp:  start,
		Request:    request.String(),
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}
