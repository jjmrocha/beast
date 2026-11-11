/*
 * Copyright 2020 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
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

import "time"

const (
	StatusGenerationError = -100
	StatusTimeout         = -400
	StatusUnexpectedError = -500
)

// Response contains the status code and the duration taken for execution of a request
type Response struct {
	Timestamp  time.Time
	Request    string
	StatusCode int
	Duration   time.Duration
}

// IsSuccess returns true for 2xx status codes
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// Failed returns true if no HTTP answer was received from the endpoint
func (r *Response) Failed() bool {
	return r.StatusCode < 0
}

// FailureReason returns the description of the failure, or an empty string when the request did not fail
func (r *Response) FailureReason() string {
	switch r.StatusCode {
	case StatusGenerationError:
		return "Request generation error"
	case StatusTimeout:
		return "Request timeout"
	case StatusUnexpectedError:
		return "Unexpected error"
	}

	return ""
}
