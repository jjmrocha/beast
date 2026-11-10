/*
 * Copyright 2026 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
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

package report

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jjmrocha/beast/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	successLine = "successfully generated"
	errorLine   = "Error writing output file"
)

func TestNewOutput(t *testing.T) {
	t.Run("returns an error when the folder does not exist", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "missing", "output.csv")
		// when
		_, err := NewOutput(fileName)
		// then
		assert.Error(t, err)
	})

	t.Run("an empty file name gives an output that writes nothing", func(t *testing.T) {
		// given
		var printed bytes.Buffer
		// when
		output, err := NewOutput("")
		// then
		require.NoError(t, err)
		output.Write(resultsOfEachKind[0].input)
		output.Close(&printed)
		assert.Empty(t, printed.String())
	})
}

func TestClose(t *testing.T) {
	t.Run("prints the success line when every write worked", func(t *testing.T) {
		// given
		var printed bytes.Buffer
		output := mustNewOutput(t, filepath.Join(t.TempDir(), "output.csv"))
		output.Write(resultsOfEachKind[0].input)
		// when
		output.Close(&printed)
		// then
		assert.Contains(t, printed.String(), successLine)
		assert.NotContains(t, printed.String(), errorLine)
	})

	t.Run("writes every record in the order it was written", func(t *testing.T) {
		// given
		const records = 5000
		var printed bytes.Buffer
		fileName := filepath.Join(t.TempDir(), "output.csv")
		output := mustNewOutput(t, fileName)
		for i := range records {
			output.Write(&client.Response{Timestamp: timestamp, Request: fmt.Sprintf("GET /%d", i), StatusCode: 200})
		}
		// when
		output.Close(&printed)
		// then
		result := mustReadCSV(t, fileName)
		require.Len(t, result, records+1)
		for i := range records {
			assert.Equal(t, fmt.Sprintf("GET /%d", i), result[i+1][1])
		}
	})

	t.Run("prints an error line when the output file could not be written", func(t *testing.T) {
		// given
		var printed bytes.Buffer
		fileName := filepath.Join(t.TempDir(), "output.csv")
		output := mustNewOutput(t, fileName)
		require.NoError(t, output.csvReport.file.Close())
		output.Write(resultsOfEachKind[0].input)
		// when
		output.Close(&printed)
		// then
		assert.Contains(t, printed.String(), errorLine+" '"+fileName+"'")
		assert.NotContains(t, printed.String(), successLine)
	})
}

func mustNewOutput(t testing.TB, fileName string) Output {
	t.Helper()
	output, err := NewOutput(fileName)
	if err != nil {
		t.Fatalf("NewOutput(%s): unexpected error: %v", fileName, err)
	}
	return output
}
