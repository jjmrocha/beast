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

package config

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	t.Run("empty file gives the defaults", func(t *testing.T) {
		// given
		fileName := "../../testdata/empty.json"
		// when
		result := Read(fileName)
		// then
		assert.Equal(t, Default(), result)
	})

	t.Run("reads every v2.6.1 key", func(t *testing.T) {
		// given
		fileName := "../../testdata/config_full.json"
		// when
		result := Read(fileName)
		// then
		assert.False(t, result.DisableCompression)
		assert.True(t, result.DisableKeepAlives)
		assert.Equal(t, 10, result.MaxConnections)
		assert.Equal(t, 5, result.MaxIdleConnections)
		assert.Equal(t, 7, result.RequestTimeout)
		assert.True(t, result.DisableCertificateCheck)
		assert.False(t, result.DisableRedirects)
	})
}

func TestGetMaxIdleConnections(t *testing.T) {
	const parallelConns = 100

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{name: "zero uses the number of parallel connections", input: 0, expected: parallelConns},
		{name: "a set value is used as it is", input: 5, expected: 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			cfg := Default()
			cfg.MaxIdleConnections = tc.input
			// when
			result := cfg.GetMaxIdleConnections(parallelConns)
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestWrite(t *testing.T) {
	t.Run("file reads back as the same configuration", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "config.json")
		// when
		Write(fileName, Default())
		// then
		result := Read(fileName)
		assert.Equal(t, Default(), result)
	})

	t.Run("file has exactly the v2.6.1 key names", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "config.json")
		// when
		Write(fileName, Default())
		// then
		content, err := os.ReadFile(fileName)
		require.NoError(t, err)
		var result map[string]any
		require.NoError(t, json.Unmarshal(content, &result))
		expected := []string{
			"disable-compression",
			"disable-keep-alives",
			"max-connections",
			"max-idle-connections",
			"request-timeout",
			"disable-certificate-check",
			"disable-redirects",
		}
		assert.ElementsMatch(t, expected, slices.Collect(maps.Keys(result)))
	})
}
