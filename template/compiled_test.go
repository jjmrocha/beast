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
	"testing"

	"github.com/jjmrocha/beast/data"
)

func TestCompileAndExecute(t *testing.T) {
	// given
	data := data.Read("../testdata/data.csv")
	request := Read("../testdata/template_post.json")
	expected := &TRequest{
		Method:   "POST",
		Endpoint: "http://someendpoint.pt/1",
		Headers: []THeader{
			{"Content-Type", "application/json"},
		},
		Body: "{\"id\": 1, \"value\": \"a1\"}",
	}
	// when
	cRequest, err := request.compile()
	if err != nil {
		t.Error(err)
	}

	result, err := cRequest.executeTemplate(1, data.Next())
	if err != nil {
		t.Error(err)
	}
	// then
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v expected %v", result, expected)
	}
}

func BenchmarkFromTemplateToClient(b *testing.B) {
	// given
	data := data.Read("../testdata/data.csv")
	request := Read("../testdata/template_post.json")
	// then
	b.ResetTimer()

	cRequest, err := request.compile()
	if err != nil {
		b.Error(err)
	}

	finalRequest, err := cRequest.executeTemplate(1, data.Next())
	if err != nil {
		b.Error(err)
	}

	_, err = finalRequest.request()
	if err != nil {
		b.Error(err)
	}
}
