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

	"gopkg.in/yaml.v3"
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
	if isJSON(fileName) {
		return readJSON(fileName)
	}

	return readYAML(fileName)
}

func readFile(fileName string) []byte {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatalf("Error reading template file %s: %v\n", fileName, err)
	}

	return data
}

// Write writes an HTTP request template to a file
func Write(fileName string, request *TRequest) {
	if isJSON(fileName) {
		writeJSON(fileName, request)
	} else {
		writeYAML(fileName, request)
	}
}

func writeFile(data []byte, fileName string) {
	if err := ioutil.WriteFile(fileName, data, 0666); err != nil {
		log.Printf("Error writing template to file %s: %v\n", fileName, err)
	}
}

// JSON

func isJSON(fileName string) bool {
	lowerCaseFileName := strings.ToLower(fileName)
	return strings.HasSuffix(lowerCaseFileName, ".json")
}

func readJSON(fileName string) *TRequest {
	data := readFile(fileName)

	var request TRequest
	if err := json.Unmarshal(data, &request); err != nil {
		log.Printf("Invalid JSON template file %s: %v\n", fileName, err)
	}

	if body, read := externalBody(request.Body); read {
		request.Body = body
	}

	return &request
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

func writeJSON(fileName string, request *TRequest) {
	data, err := json.MarshalIndent(request, "", "\t")
	if err != nil {
		log.Printf("Error encoding request %v to JSON: %v\n", request, err)
	}

	writeFile(data, fileName)
}

// YAML

type yamlRequest struct {
	Method   string            `yaml:"method"`
	Endpoint string            `yaml:"endpoint"`
	Headers  map[string]string `yaml:"headers,omitempty"`
	Body     string            `yaml:"request-body,omitempty"`
}

func toYamlRequest(request *TRequest) *yamlRequest {
	return &yamlRequest{
		Method:   request.Method,
		Endpoint: request.Endpoint,
		Headers:  toHeaderMap(request.Headers),
		Body:     request.Body,
	}
}

func toHeaderMap(headers []THeader) map[string]string {
	if headers == nil {
		return nil
	}

	headerMap := make(map[string]string)

	for _, header := range headers {
		headerMap[header.Key] = header.Value
	}

	return headerMap
}

func fromYamlRequest(request *yamlRequest) *TRequest {
	return &TRequest{
		Method:   request.Method,
		Endpoint: request.Endpoint,
		Headers:  fromHeaderMap(request.Headers),
		Body:     request.Body,
	}
}

func fromHeaderMap(headers map[string]string) []THeader {
	if headers == nil {
		return nil
	}

	tHeaders := make([]THeader, 0, len(headers))

	for key, value := range headers {
		header := THeader{
			Key:   key,
			Value: value,
		}
		tHeaders = append(tHeaders, header)
	}

	return tHeaders
}

func writeYAML(fileName string, request *TRequest) {
	yamlRequest := toYamlRequest(request)

	data, err := yaml.Marshal(yamlRequest)
	if err != nil {
		log.Printf("Error encoding request %v to YAML: %v\n", request, err)
	}

	writeFile(data, fileName)
}

func readYAML(fileName string) *TRequest {
	data := readFile(fileName)

	var request yamlRequest
	if err := yaml.Unmarshal(data, &request); err != nil {
		log.Printf("Invalid YAML template file %s: %v\n", fileName, err)
	}

	return fromYamlRequest(&request)
}
