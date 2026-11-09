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

	"github.com/jjmrocha/beast/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	// given
	gnt := &Generator{
		RecordID: 1,
		Data:     data.Read("../testdata/simple.csv").Next(),
	}
	// when
	result := gnt.Log()
	// then
	expected := "requestId: 1 and data: {A: a1}"
	assert.Equal(t, expected, result)
}

func TestRequest(t *testing.T) {
	t.Run("builds the request from the rendered template", func(t *testing.T) {
		// given
		gnt := &Generator{
			Template: mustCompile(t, "../testdata/template_post.json"),
			RecordID: 1,
			Data:     data.Read("../testdata/data.csv").Next(),
		}
		// when
		result, err := gnt.Request()
		// then
		require.NoError(t, err)
		expected := "POST http://someendpoint.pt/1"
		assert.Equal(t, expected, result.String())
	})

	t.Run("returns an error when the rendered URL is invalid", func(t *testing.T) {
		// given
		tmpl := &Template{Method: "GET", Endpoint: "http://[::1"}
		tmplc, err := tmpl.Compile()
		require.NoError(t, err)
		record := data.NewRecord()
		gnt := &Generator{Template: tmplc, RecordID: 1, Data: &record}
		// when
		result, err := gnt.Request()
		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func BenchmarkRequest(b *testing.B) {
	gnt := &Generator{
		Template: mustCompile(b, "../testdata/template_post.json"),
		RecordID: 1,
		Data:     data.Read("../testdata/data.csv").Next(),
	}

	for b.Loop() {
		if _, err := gnt.Request(); err != nil {
			b.Fatal(err)
		}
	}
}
