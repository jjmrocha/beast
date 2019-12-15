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

// Package request provide functions to read, write templates and generate requests
package request

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

// THeader represents an HTTP Header template
type THeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TRequest represents an HTTP request template
type TRequest struct {
	Method   string    `json:"method"`
	Endpoint string    `json:"url"`
	Headers  []THeader `json:"headers,omitempty"`
	Body     string    `json:"body,omitempty"`
}

// Read reads an HTTP request template from a file
func Read(fileName string) *TRequest {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading template file %s: %v\n", fileName, err)
	}

	var request TRequest
	json.Unmarshal(data, &request)

	if body, readed := externalBody(request.Body); readed {
		request.Body = body
	}

	return &request
}

// Write writes an HTTP request template to a file
func Write(fileName string, request *TRequest) {
	data, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		log.Printf("Error encoding request %v to JSON: %v\n", request, err)
	}

	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		log.Printf("Error writing template to file %s: %v\n", fileName, err)
	}
}

func externalBody(body string) (string, bool) {
	if strings.HasPrefix(body, "@") {
		fileName := body[1:]
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatalf("Error reading external body file %s: %v\n", fileName, err)
		}

		return string(data), true
	}

	return "", false
}
