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

package client

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/jjmrocha/beast/config"
)

func TestHTTPTimeout(t *testing.T) {
	// given
	var tests = []struct {
		input    int
		expected time.Duration
	}{
		{1, 1 * time.Second},
		{30, 30 * time.Second},
	}
	// then
	for _, test := range tests {
		config := config.Default()
		config. = test.input
		result := HTTP(config)
		if client, ok := result.native.(*http.Client); ok && client.Timeout != test.expected {
			t.Errorf("got %v expected %v for RequestTimeout", client.Timeout, test.expected)
		}
	}
}

// Mocked error for generation of timeouts
type timeoutMockedError bool

func (t timeoutMockedError) Timeout() bool {
	return bool(t)
}

func (t timeoutMockedError) Error() string {
	return fmt.Sprintf("timeout = %v", bool(t))
}

// Mocked httpClient
type timeoutMockedClient bool

func (t timeoutMockedClient) Do(r *http.Request) (*http.Response, error) {
	err := timeoutMockedError(t)
	return nil, &url.Error{Err: &err}
}

func TestExecuteTimeout(t *testing.T) {
	// given
	client := &BClient{native: timeoutMockedClient(true)}
	expected := -400
	// when
	result := client.Execute(&BRequest{})
	// then
	if result.StatusCode != expected {
		t.Errorf("got %v expected %v", result.StatusCode, expected)
	}
}

func TestExecuteGeneric(t *testing.T) {
	// given
	client := &BClient{native: timeoutMockedClient(false)}
	expected := -500
	// when
	result := client.Execute(&BRequest{})
	// then
	if result.StatusCode != expected {
		t.Errorf("got %v expected %v", result.StatusCode, expected)
	}
}
