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

package template

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadBasicGETforJSON(t *testing.T) {
	// given
	expected := &Template{
		Method:   "GET",
		Endpoint: "http://www.google.pt",
	}
	// when
	tmpl := Read("../testdata/basic_get.json")
	// then
	if !reflect.DeepEqual(tmpl, expected) {
		t.Errorf("got %v expected %v", tmpl, expected)
	}
}

func TestReadBasicGETforYAML(t *testing.T) {
	// given
	expected := &Template{
		Method:   "GET",
		Endpoint: "http://www.google.pt",
	}
	// when
	tmpl := Read("../testdata/basic_get.yaml")
	// then
	if !reflect.DeepEqual(tmpl, expected) {
		t.Errorf("got %v expected %v", tmpl, expected)
	}
}

func TestReadBasicPOSTforJSON(t *testing.T) {
	// given
	expected := &Template{
		Method:   "POST",
		Endpoint: "http://someendpoint.pt",
		Headers: []Header{
			{"Content-Type", "application/json"},
		},
		Body: "{\"id\": 1, \"value\": \"any\"}",
	}
	// when
	tmpl := Read("../testdata/basic_post.json")
	// then
	if !reflect.DeepEqual(tmpl, expected) {
		t.Errorf("got %v expected %v", tmpl, expected)
	}
}

func TestReadBasicPOSTforYAML(t *testing.T) {
	// given
	expected := &Template{
		Method:   "POST",
		Endpoint: "http://someendpoint.pt",
		Headers: []Header{
			{"Content-Type", "application/json"},
		},
		Body: "{\"id\": 1, \"value\": \"any\"}",
	}
	// when
	tmpl := Read("../testdata/basic_post.yaml")
	// then
	if !reflect.DeepEqual(tmpl, expected) {
		t.Errorf("got %v expected %v", tmpl, expected)
	}
}

func TestReadTemplatePOSTforJSON(t *testing.T) {
	// given
	expected := &Template{
		Method:   "POST",
		Endpoint: "http://someendpoint.pt/{{ .RequestID }}",
		Headers: []Header{
			{"Content-Type", "application/json"},
		},
		Body: "{\"id\": {{ .RequestID }}, \"value\": \"{{ .Data.A }}\"}",
	}
	// when
	tmpl := Read("../testdata/template_post.json")
	// then
	if !reflect.DeepEqual(tmpl, expected) {
		t.Errorf("got %v expected %v", tmpl, expected)
	}
}

func TestReadTemplatePOSTforYAML(t *testing.T) {
	// given
	expectedMethod := "POST"
	expectedEndpoint := "http://someendpoint.pt/{{ .RequestID }}"
	expectedHeaders := []Header{
		{"Content-Type", "application/json"},
	}
	firstBodyField := "\"id\": {{ .RequestID }}"
	secondBodyField := "\"value\": \"{{ .Data.A }}\""
	// when
	tmpl := Read("../testdata/template_post.yaml")
	// then
	if tmpl.Method != expectedMethod {
		t.Errorf("got %v expected %v", tmpl.Method, expectedMethod)
	}

	if tmpl.Endpoint != expectedEndpoint {
		t.Errorf("got %v expected %v", tmpl.Endpoint, expectedEndpoint)
	}

	if !reflect.DeepEqual(tmpl.Headers, expectedHeaders) {
		t.Errorf("got %v expected %v", tmpl.Headers, expectedHeaders)
	}

	if !strings.Contains(tmpl.Body, firstBodyField) {
		t.Errorf("body %v should contain %v", tmpl.Body, firstBodyField)
	}

	if !strings.Contains(tmpl.Body, secondBodyField) {
		t.Errorf("body %v should contain %v", tmpl.Body, secondBodyField)
	}
}

func TestReadTemplatePOSTWithExternalBody(t *testing.T) {
	// given
	expected := &Template{
		Method:   "POST",
		Endpoint: "http://someendpoint.pt/{{ .RequestID }}",
		Headers: []Header{
			{"Content-Type", "application/json"},
		},
		Body: "{\"id\": {{ .RequestID }}, \"value\": \"{{ .Data.A }}\"}",
	}
	// when
	tmpl := Read("../testdata/template_post_external.json")
	// then
	if !reflect.DeepEqual(tmpl, expected) {
		t.Errorf("got %v expected %v", tmpl, expected)
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
	tHeader := []Header{
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
	tRequest := &Template{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
		Headers: []Header{
			{Key: "User-Agent", Value: "Beast/1"},
		},
		Body: "Optional, enter body to send with POST or PUT",
	}
	yRequest := &templateY{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
		Headers: map[string]string{
			"User-Agent": "Beast/1",
		},
		Body: "Optional, enter body to send with POST or PUT",
	}
	// when
	yResult := toYamlTemplate(tRequest)
	tResult := fromYamlTemplate(yRequest)
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
	tRequest := &Template{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
	}
	yRequest := &templateY{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
	}
	// then
	yResult := toYamlTemplate(tRequest)
	tResult := fromYamlTemplate(yRequest)
	// then
	if !reflect.DeepEqual(yResult, yRequest) {
		t.Errorf("got %v expected %v", yResult, yRequest)
	}

	if !reflect.DeepEqual(tResult, tRequest) {
		t.Errorf("got %v expected %v", tResult, tRequest)
	}
}
