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
	"testing"

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/data"
)

func TestLog(t *testing.T) {
	// given
	dt := data.Read("../testdata/data.csv")
	gnt := &Generator{
		recordID: 1,
		data:     dt.Next(),
	}
	expected := "requestId: 1 and data: {A: a1, B: b1}"
	// when
	result := gnt.Log()
	// then
	if result != expected {
		t.Errorf("got %v expected %v", result, expected)
	}
}

func TestRequestForStatic(t *testing.T) {
	// given
	expected := &client.Request{}
	gnt := &Generator{
		final: expected,
	}
	// when
	result, err := gnt.Request()
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
	dt := data.Read("../testdata/data.csv")
	tmpl := Read("../testdata/template_post.json")
	tmplc, _ := tmpl.compile()
	gnt := &Generator{
		data:     dt.Next(),
		recordID: 1,
		template: tmplc,
	}
	expected := "POST http://someendpoint.pt/1"
	// when
	result, err := gnt.Request()
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
	dt := data.Read("../testdata/data.csv")
	tmpl := Read("../testdata/template_post.json")
	tmplc, _ := tmpl.compile()
	gnt := &Generator{
		data:     dt.Next(),
		recordID: 1,
		template: tmplc,
	}
	// when
	b.ResetTimer()
	_, err := gnt.Request()
	// then
	if err != nil {
		b.Error(err)
	}
}
