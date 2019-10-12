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
	"testing"
)

func TestReadBasicGET(t *testing.T) {
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

func TestReadBasicPOST(t *testing.T) {
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

func TestReadTemplatePOST(t *testing.T) {
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
