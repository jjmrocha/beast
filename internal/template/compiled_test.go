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

	"github.com/jjmrocha/beast/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteTemplate(t *testing.T) {
	t.Run("fills the URL, headers and body", func(t *testing.T) {
		// given
		tmplc := mustCompile(t, "../../testdata/template_post.json")
		record := data.Read("../../testdata/data.csv").Next()
		// when
		result, err := tmplc.executeTemplate(1, record)
		// then
		require.NoError(t, err)
		expected := &Template{
			Method:   "POST",
			Endpoint: "http://someendpoint.pt/1",
			Headers:  jsonContentType,
			Body:     "{\"id\": 1, \"value\": \"a1\"}",
		}
		assert.Equal(t, expected, result)
	})

	t.Run("keeps an empty body when the template has none", func(t *testing.T) {
		// given
		tmplc := mustCompile(t, "../../testdata/basic_get.json")
		record := data.NewRecord()
		// when
		result, err := tmplc.executeTemplate(1, &record)
		// then
		require.NoError(t, err)
		assert.Empty(t, result.Body)
	})

	t.Run("returns an error when the template fails while running", func(t *testing.T) {
		// given
		tmpl := &Template{Method: "GET", Endpoint: "http://someendpoint.pt/{{ .Data.A.B }}"}
		tmplc, err := tmpl.Compile()
		require.NoError(t, err)
		record := data.Record{"A": "a1"}
		// when
		result, err := tmplc.executeTemplate(1, &record)
		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func BenchmarkFromTemplateToClient(b *testing.B) {
	tmpl := Read("../../testdata/template_post.json")
	record := data.Read("../../testdata/data.csv").Next()

	for b.Loop() {
		tmplc, err := tmpl.Compile()
		if err != nil {
			b.Fatal(err)
		}
		tmplf, err := tmplc.executeTemplate(1, record)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := tmplf.request(); err != nil {
			b.Fatal(err)
		}
	}
}

func mustCompile(t testing.TB, fileName string) *CompiledTemplate {
	t.Helper()
	tmplc, err := Read(fileName).Compile()
	if err != nil {
		t.Fatalf("Compile(%s): unexpected error: %v", fileName, err)
	}
	return tmplc
}
