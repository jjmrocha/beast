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
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jjmrocha/beast/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	formattedTimestamp = "2026-11-09 10:30:15.123"
	getRequest         = "GET http://someendpoint.pt/1"
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
		expected: []string{formattedTimestamp, getRequest, "Executed", "404", "false", "20"},
	},
	{
		name:     "request generation error",
		input:    &client.Response{Timestamp: timestamp, StatusCode: -100},
		expected: []string{formattedTimestamp, "", "Request generation error", "", "false", ""},
	},
	{
		name:     "request timeout",
		input:    &client.Response{Timestamp: timestamp, Request: getRequest, StatusCode: -400, Duration: time.Second},
		expected: []string{formattedTimestamp, getRequest, "Request timeout", "", "false", "1000"},
	},
	{
		name:     "unexpected error",
		input:    &client.Response{Timestamp: timestamp, Request: getRequest, StatusCode: -500, Duration: 5 * time.Millisecond},
		expected: []string{formattedTimestamp, getRequest, "Unexpected error", "", "false", "5"},
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
		report := newCsvReport(fileName)
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

	t.Run("quotes a request that contains a comma", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "output.csv")
		report := newCsvReport(fileName)
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
