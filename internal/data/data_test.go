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
		result, err := Read(fileName)
		// then
		require.NoError(t, err)
		assert.Equal(t, &firstRow, result.Next())
		assert.Equal(t, &secondRow, result.Next())
	})

	t.Run("keeps a quoted value with a comma as one value", func(t *testing.T) {
		// given
		fileName := writeCSV(t, "NAME,URL\n\"Smith, John\",http://a.pt\n")
		// when
		result, err := Read(fileName)
		// then
		require.NoError(t, err)
		expected := &Record{"NAME": "Smith, John", "URL": "http://a.pt"}
		assert.Equal(t, expected, result.Next())
	})

	t.Run("a file with only a header returns ErrNoRows", func(t *testing.T) {
		// given
		fileName := writeCSV(t, "A,B\n")
		// when
		result, err := Read(fileName)
		// then
		assert.ErrorIs(t, err, ErrNoRows)
		assert.Nil(t, result)
	})

	t.Run("an empty file returns ErrNoRows", func(t *testing.T) {
		// given
		fileName := writeCSV(t, "")
		// when
		result, err := Read(fileName)
		// then
		assert.ErrorIs(t, err, ErrNoRows)
		assert.Nil(t, result)
	})

	t.Run("a missing file returns an error", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "missing.csv")
		// when
		result, err := Read(fileName)
		// then
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, result)
	})
}

func TestNext(t *testing.T) {
	t.Run("starts again at the first row after the last", func(t *testing.T) {
		// given
		dt := mustRead(t, "../../testdata/data.csv")
		dt.Next()
		dt.Next()
		// when
		result := dt.Next()
		// then
		assert.Equal(t, &firstRow, result)
	})

	t.Run("returns the only row every time", func(t *testing.T) {
		// given
		dt := mustRead(t, writeCSV(t, "A\na1\n"))
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

func mustRead(t testing.TB, fileName string) *Data {
	t.Helper()
	dt, err := Read(fileName)
	if err != nil {
		t.Fatalf("Read(%s): unexpected error: %v", fileName, err)
	}
	return dt
}
