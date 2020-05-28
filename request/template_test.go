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

package request

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadBasicGETforJSON(t *testing.T) {
	// given
	expected := &TRequest{
		Method:   "GET",
		Endpoint: "http://www.google.pt",
	}
	// when
	template := Read("../testdata/basic_get.json")
	// then
	if !reflect.DeepEqual(template, expected) {
		t.Errorf("got %v expected %v", template, expected)
	}
}

func TestReadBasicGETforYAML(t *testing.T) {
	// given
	expected := &TRequest{
		Method:   "GET",
		Endpoint: "http://www.google.pt",
	}
	// when
	template := Read("../testdata/basic_get.yaml")
	// then
	if !reflect.DeepEqual(template, expected) {
		t.Errorf("got %v expected %v", template, expected)
	}
}

func TestReadBasicPOSTforJSON(t *testing.T) {
	// given
	expected := &TRequest{
		Method:   "POST",
		Endpoint: "http://someendpoint.pt",
		Headers: []THeader{
			{"Content-Type", "application/json"},
		},
		Body: "{\"id\": 1, \"value\": \"any\"}",
	}
	// when
	template := Read("../testdata/basic_post.json")
	// then
	if !reflect.DeepEqual(template, expected) {
		t.Errorf("got %v expected %v", template, expected)
	}
}

func TestReadBasicPOSTforYAML(t *testing.T) {
	// given
	expected := &TRequest{
		Method:   "POST",
		Endpoint: "http://someendpoint.pt",
		Headers: []THeader{
			{"Content-Type", "application/json"},
		},
		Body: "{\"id\": 1, \"value\": \"any\"}",
	}
	// when
	template := Read("../testdata/basic_post.yaml")
	// then
	if !reflect.DeepEqual(template, expected) {
		t.Errorf("got %v expected %v", template, expected)
	}
}

func TestReadTemplatePOSTforJSON(t *testing.T) {
	// given
	expected := &TRequest{
		Method:   "POST",
		Endpoint: "http://someendpoint.pt/{{ .RequestID }}",
		Headers: []THeader{
			{"Content-Type", "application/json"},
		},
		Body: "{\"id\": {{ .RequestID }}, \"value\": \"{{ .Data.A }}\"}",
	}
	// when
	template := Read("../testdata/template_post.json")
	// then
	if !reflect.DeepEqual(template, expected) {
		t.Errorf("got %v expected %v", template, expected)
	}
}

func TestReadTemplatePOSTforYAML(t *testing.T) {
	// given
	expectedMethod := "POST"
	expectedEndpoint := "http://someendpoint.pt/{{ .RequestID }}"
	expectedHeaders := []THeader{
		{"Content-Type", "application/json"},
	}
	firstBodyField := "\"id\": {{ .RequestID }}"
	secondBodyField := "\"value\": \"{{ .Data.A }}\""
	// when
	template := Read("../testdata/template_post.yaml")
	// then
	if template.Method != expectedMethod {
		t.Errorf("got %v expected %v", template.Method, expectedMethod)
	}

	if template.Endpoint != expectedEndpoint {
		t.Errorf("got %v expected %v", template.Endpoint, expectedEndpoint)
	}

	if !reflect.DeepEqual(template.Headers, expectedHeaders) {
		t.Errorf("got %v expected %v", template.Headers, expectedHeaders)
	}

	if !strings.Contains(template.Body, firstBodyField) {
		t.Errorf("body %v should contain %v", template.Body, firstBodyField)
	}

	if !strings.Contains(template.Body, secondBodyField) {
		t.Errorf("body %v should contain %v", template.Body, secondBodyField)
	}
}

func TestReadTemplatePOSTWithExternalBody(t *testing.T) {
	// given
	expected := &TRequest{
		Method:   "POST",
		Endpoint: "http://someendpoint.pt/{{ .RequestID }}",
		Headers: []THeader{
			{"Content-Type", "application/json"},
		},
		Body: "{\"id\": {{ .RequestID }}, \"value\": \"{{ .Data.A }}\"}",
	}
	// when
	template := Read("../testdata/template_post_external.json")
	// then
	if !reflect.DeepEqual(template, expected) {
		t.Errorf("got %v expected %v", template, expected)
	}
}

func TestExternalBodyWithExternalBody(t *testing.T) {
	// given
	body := "@../testdata/body.json"
	expected := "{\"id\": {{ .RequestID }}, \"value\": \"{{ .Data.A }}\"}"
	// when
	response, found := externalBody(body)
	// then
	if found != true {
		t.Errorf("got %v expected %v for found", found, true)
	}

	if response != expected {
		t.Errorf("got %v expected %v for response", response, expected)
	}
}

func TestExternalBodyWithoutExternalBody(t *testing.T) {
	// given
	body := "{\"id\": {{ .RequestID }}, \"value\": \"{{ .Data.A }}\"}"
	// when
	_, found := externalBody(body)
	// then
	if found != false {
		t.Errorf("got %v expected %v", found, false)
	}
}

func TestIsJSON(t *testing.T) {
	// given
	var tests = []struct {
		input    string
		expected bool
	}{
		{"../testdata/basic_get.json", true},
		{"../testdata/basic_get.yaml", false},
		{"template_post_external.json", true},
		{"template_post_external.yaml", false},
		{"SOME_FILE.JSON", true},
		{"/json/not_ajson", false},
	}
	// then
	for _, test := range tests {
		result := isJSON(test.input)
		if result != test.expected {
			t.Errorf("got %v expected %v", result, test.expected)
		}
	}
}

func TestHeaderConvertion(t *testing.T) {
	// given
	tHeader := []THeader{
		{
			Key:   "Key_1",
			Value: "Value_1",
		},
		{
			Key:   "Key_2",
			Value: "Value_2",
		},
	}
	yHeader := map[string]string{
		"Key_1": "Value_1",
		"Key_2": "Value_2",
	}
	// when
	yResult := toHeaderMap(tHeader)
	tResult := fromHeaderMap(yHeader)
	// then
	if !reflect.DeepEqual(yResult, yHeader) {
		t.Errorf("got %v expected %v", yResult, yHeader)
	}

	if !reflect.DeepEqual(tResult, tHeader) {
		t.Errorf("got %v expected %v", tResult, tHeader)
	}
}

func TestRequestConvertion(t *testing.T) {
	// given
	tRequest := &TRequest{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
		Headers: []THeader{
			{Key: "User-Agent", Value: "Beast/1"},
		},
		Body: "Optional, enter body to send with POST or PUT",
	}
	yRequest := &yamlRequest{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
		Headers: map[string]string{
			"User-Agent": "Beast/1",
		},
		Body: "Optional, enter body to send with POST or PUT",
	}
	// when
	yResult := toYamlRequest(tRequest)
	tResult := fromYamlRequest(yRequest)
	// then
	if !reflect.DeepEqual(yResult, yRequest) {
		t.Errorf("got %v expected %v", yResult, yRequest)
	}

	if !reflect.DeepEqual(tResult, tRequest) {
		t.Errorf("got %v expected %v", tResult, tRequest)
	}
}

func TestPartialRequestConvertion(t *testing.T) {
	// given
	tRequest := &TRequest{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
	}
	yRequest := &yamlRequest{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
	}
	// then
	yResult := toYamlRequest(tRequest)
	tResult := fromYamlRequest(yRequest)
	// then
	if !reflect.DeepEqual(yResult, yRequest) {
		t.Errorf("got %v expected %v", yResult, yRequest)
	}

	if !reflect.DeepEqual(tResult, tRequest) {
		t.Errorf("got %v expected %v", tResult, tRequest)
	}
}
