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

package client

import (
	"fmt"
	"time"
)

// Response contains the status code and the duration taken for execution of a request
type Response struct {
	Timestamp  time.Time
	Request    string
	StatusCode int
	Duration   time.Duration
}

func (r *Response) String() string {
	return fmt.Sprintf("%v - %v", r.StatusCode, r.Duration)
}

// IsSuccess return true for statusCodes matchs 2xx
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError returns true if we didn't receive an awnser from the endpoint
func (r *Response) IsClientError() bool {
	return r.StatusCode < 0
}

// ClientError returns the descriprion for the client error
func (r *Response) ClientError() string {
	switch r.StatusCode {
	case -100:
		return "Request generation error"
	case -400:
		return "Request timeout"
	case -500:
		return "Unexpected error"
	}

	return ""
}
