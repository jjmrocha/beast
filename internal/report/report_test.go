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
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jjmrocha/beast/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	formattedTimestamp = "2026-11-09 10:30:15.123"
	getRequest         = "GET http://someendpoint.pt/1"
	notSuccess         = "false"
)

var timestamp = time.Date(2026, 11, 9, 10, 30, 15, 123_000_000, time.UTC)

var resultsOfEachKind = []struct {
	name     string
	input    *client.Response
	expected []string
}{
	{
		name:     "success",
		input:    &client.Response{Timestamp: timestamp, Request: getRequest, StatusCode: 200, Duration: 1500 * time.Millisecond},
		expected: []string{formattedTimestamp, getRequest, "Executed", "200", "true", "1500"},
	},
	{
		name:     "non-success status",
		input:    &client.Response{Timestamp: timestamp, Request: getRequest, StatusCode: 404, Duration: 20 * time.Millisecond},
		expected: []string{formattedTimestamp, getRequest, "Executed", "404", notSuccess, "20"},
	},
	{
		name:     "request generation error",
		input:    &client.Response{Timestamp: timestamp, StatusCode: client.StatusGenerationError},
		expected: []string{formattedTimestamp, "", "Request generation error", "", notSuccess, ""},
	},
	{
		name:     "request timeout",
		input:    &client.Response{Timestamp: timestamp, Request: getRequest, StatusCode: client.StatusTimeout, Duration: time.Second},
		expected: []string{formattedTimestamp, getRequest, "Request timeout", "", notSuccess, "1000"},
	},
	{
		name:     "unexpected error",
		input:    &client.Response{Timestamp: timestamp, Request: getRequest, StatusCode: client.StatusUnexpectedError, Duration: 5 * time.Millisecond},
		expected: []string{formattedTimestamp, getRequest, "Unexpected error", "", notSuccess, "5"},
	},
}

func TestEncode(t *testing.T) {
	for _, tc := range resultsOfEachKind {
		t.Run(tc.name, func(t *testing.T) {
			// given
			response := tc.input
			// when
			result := encode(response)
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCsvReport(t *testing.T) {
	t.Run("writes the v2.6.1 header and one line per result", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "output.csv")
		report := mustNewCsvReport(t, fileName)
		// when
		for _, tc := range resultsOfEachKind {
			report.write(tc.input)
		}
		report.flush()
		// then
		result := mustReadCSV(t, fileName)
		expected := [][]string{{"Timestamp", "Request", "Result", "StatusCode", "IsSuccess", "Duration"}}
		for _, tc := range resultsOfEachKind {
			expected = append(expected, tc.expected)
		}
		assert.Equal(t, expected, result)
	})

	t.Run("keeps the error when the file cannot be written", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "output.csv")
		report := mustNewCsvReport(t, fileName)
		require.NoError(t, report.file.Close())
		report.write(resultsOfEachKind[0].input)
		// when
		report.flush()
		// then
		assert.Error(t, report.err)
	})

	t.Run("has no error when the file was written", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "output.csv")
		report := mustNewCsvReport(t, fileName)
		report.write(resultsOfEachKind[0].input)
		// when
		report.flush()
		// then
		assert.NoError(t, report.err)
	})

	t.Run("quotes a request that contains a comma", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "output.csv")
		report := mustNewCsvReport(t, fileName)
		request := "GET http://someendpoint.pt/a,b"
		response := &client.Response{Timestamp: timestamp, Request: request, StatusCode: 200, Duration: time.Millisecond}
		// when
		report.write(response)
		report.flush()
		// then
		content, err := os.ReadFile(fileName)
		require.NoError(t, err)
		assert.Contains(t, string(content), `"GET http://someendpoint.pt/a,b"`)
		result := mustReadCSV(t, fileName)
		require.Len(t, result, 2)
		assert.Equal(t, request, result[1][1])
	})
}

func mustReadCSV(t testing.TB, fileName string) [][]string {
	t.Helper()
	content, err := os.ReadFile(fileName)
	require.NoError(t, err)
	records, err := csv.NewReader(strings.NewReader(string(content))).ReadAll()
	require.NoError(t, err)
	return records
}

func mustNewCsvReport(t testing.TB, fileName string) *csvReport {
	t.Helper()
	report, err := newCsvReport(fileName)
	if err != nil {
		t.Fatalf("newCsvReport(%s): unexpected error: %v", fileName, err)
	}
	return report
}
