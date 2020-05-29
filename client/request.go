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
	"net/http"
)

// BRequest represents an HTTP request
type BRequest struct {
	native *http.Request
}

// MakeRequest creates a BRequest using a http.Request
func MakeRequest(req *http.Request) *BRequest {
	return &BRequest{native: req}
}

func (r *BRequest) String() string {
	return fmt.Sprintf("%s %s", r.native.Method, r.native.URL)
}
