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
	"testing"

	"github.com/jjmrocha/beast/client"

	"github.com/jjmrocha/beast/data"
)

func TestLog(t *testing.T) {
	// given
	data := data.Read("../testdata/data.csv")
	generator := &Generator{
		recordID: 1,
		data:     data.Next(),
	}
	expected := "requestId: 1 and data: {A: a1, B: b1}"
	// when
	result := generator.Log()
	// then
	if result != expected {
		t.Errorf("got %v expected %v", result, expected)
	}
}

func TestRequestForStatic(t *testing.T) {
	// given
	expected := &client.BRequest{}
	generator := &Generator{
		final: expected,
	}
	// when
	result, err := generator.Request()
	// then
	if err != nil {
		t.Errorf("Error not expected: %v", err)
	}

	if result != expected {
		t.Errorf("got %v expected %v", result, expected)
	}
}

func TestRequestForDynamic(t *testing.T) {
	// given
	data := data.Read("../testdata/data.csv")
	request := Read("../testdata/template_post.json")
	cRequest, _ := request.compile()
	generator := &Generator{
		data:     data.Next(),
		recordID: 1,
		template: cRequest,
	}
	expected := "POST http://someendpoint.pt/1"
	// when
	result, err := generator.Request()
	// then
	if err != nil {
		t.Errorf("Error not expected: %v", err)
	}

	if result.String() != expected {
		t.Errorf("got %v expected %v", result, expected)
	}
}

func BenchmarkRequest(b *testing.B) {
	// given
	data := data.Read("../testdata/data.csv")
	request := Read("../testdata/template_post.json")
	cRequest, _ := request.compile()
	generator := &Generator{
		data:     data.Next(),
		recordID: 1,
		template: cRequest,
	}
	// when
	b.ResetTimer()
	_, err := generator.Request()
	// then
	if err != nil {
		b.Error(err)
	}
}
