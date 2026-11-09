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

package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	firstRow  = Record{"A": "a1", "B": "b1"}
	secondRow = Record{"A": "a2", "B": "b2"}
)

func TestRead(t *testing.T) {
	t.Run("uses the first line as column names and returns the rows in file order", func(t *testing.T) {
		// given
		fileName := "../../testdata/data.csv"
		// when
		result := Read(fileName)
		// then
		assert.Equal(t, &firstRow, result.Next())
		assert.Equal(t, &secondRow, result.Next())
	})

	t.Run("keeps a quoted value with a comma as one value", func(t *testing.T) {
		// given
		fileName := writeCSV(t, "NAME,URL\n\"Smith, John\",http://a.pt\n")
		// when
		result := Read(fileName)
		// then
		expected := &Record{"NAME": "Smith, John", "URL": "http://a.pt"}
		assert.Equal(t, expected, result.Next())
	})
}

func TestNext(t *testing.T) {
	t.Run("starts again at the first row after the last", func(t *testing.T) {
		// given
		dt := Read("../../testdata/data.csv")
		dt.Next()
		dt.Next()
		// when
		result := dt.Next()
		// then
		assert.Equal(t, &firstRow, result)
	})

	t.Run("returns the only row every time", func(t *testing.T) {
		// given
		dt := Read(writeCSV(t, "A\na1\n"))
		dt.Next()
		// when
		result := dt.Next()
		// then
		expected := &Record{"A": "a1"}
		assert.Equal(t, expected, result)
	})
}

func writeCSV(t testing.TB, content string) string {
	t.Helper()
	fileName := filepath.Join(t.TempDir(), "data.csv")
	require.NoError(t, os.WriteFile(fileName, []byte(content), 0o600))
	return fileName
}
