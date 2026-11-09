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

package report

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	successLine = "successfully generated"
	errorLine   = "Error writing output file"
)

func TestClose(t *testing.T) {
	t.Run("prints the success line when every write worked", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "output.csv")
		output := NewOutput(fileName)
		output.Write(resultsOfEachKind[0].input)
		// when
		result := captureStdout(t, output.Close)
		// then
		assert.Contains(t, result, successLine)
		assert.NotContains(t, result, errorLine)
	})

	t.Run("prints an error line when the output file could not be written", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "output.csv")
		output := NewOutput(fileName)
		require.NoError(t, output.csvReport.file.Close())
		output.Write(resultsOfEachKind[0].input)
		// when
		result := captureStdout(t, output.Close)
		// then
		assert.Contains(t, result, errorLine+" '"+fileName+"'")
		assert.NotContains(t, result, successLine)
	})

	t.Run("prints nothing when there is no output file", func(t *testing.T) {
		// given
		output := NewOutput("")
		// when
		result := captureStdout(t, output.Close)
		// then
		assert.Empty(t, result)
	})
}

func captureStdout(t testing.TB, fn func()) string {
	t.Helper()
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	original := os.Stdout
	os.Stdout = writer
	t.Cleanup(func() { os.Stdout = original })

	done := make(chan string)
	go func() {
		data, _ := io.ReadAll(reader)
		done <- string(data)
	}()

	fn()
	os.Stdout = original
	require.NoError(t, writer.Close())
	return <-done
}
