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
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jjmrocha/beast/internal/config"
)

// Client represents an HTTP client
type Client struct {
	native *http.Client
}

// NewClient creates a client.Client based on the provided configuration
func NewClient(cfg *config.Config, parallelConns int) *Client {
	tls := &tls.Config{
		InsecureSkipVerify: cfg.DisableCertificateCheck, //nolint:gosec // only when the user sets disable-certificate-check
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
func (c *Client) Execute(request *http.Request) *Response {
	description := request.Method + " " + request.URL.String()
	start := time.Now()
	resp, err := c.native.Do(request)
	duration := time.Since(start)

	if err != nil {
		var urlErr *url.Error
		if ok := errors.As(err, &urlErr); ok && urlErr.Timeout() {
			return &Response{
				Timestamp:  start,
				Request:    description,
				StatusCode: StatusTimeout,
				Duration:   duration,
			}
		}

		log.Printf("Error executing request '%s': %v\n", description, err) //nolint:gosec // the URL comes from the user's own template
		return &Response{
			Timestamp:  start,
			Request:    description,
			StatusCode: StatusUnexpectedError,
			Duration:   duration,
		}
	}

	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)

	return &Response{
		Timestamp:  start,
		Request:    description,
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}
