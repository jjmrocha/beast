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

package config

import (
	"bytes"
	"encoding/json"
	"log"
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
		result, err := Read(fileName)
		// then
		require.NoError(t, err)
		assert.Equal(t, Default(), result)
	})

	t.Run("logs invalid JSON and keeps the defaults", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "config.json")
		require.NoError(t, os.WriteFile(fileName, []byte("{"), 0o600))
		var logged bytes.Buffer
		original := log.Writer()
		log.SetOutput(&logged)
		t.Cleanup(func() { log.SetOutput(original) })
		// when
		result, err := Read(fileName)
		// then
		require.NoError(t, err)
		assert.Equal(t, Default(), result)
		assert.Contains(t, logged.String(), "Invalid JSON config file "+fileName)
	})

	t.Run("a missing file returns an error", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "missing.json")
		// when
		result, err := Read(fileName)
		// then
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, result)
	})

	for _, key := range []string{"max-connections", "max-idle-connections", "request-timeout"} {
		t.Run("a negative "+key+" returns ErrInvalidConfig", func(t *testing.T) {
			// given
			fileName := filepath.Join(t.TempDir(), "config.json")
			require.NoError(t, os.WriteFile(fileName, []byte(`{"`+key+`": -1}`), 0o600))
			// when
			result, err := Read(fileName)
			// then
			assert.ErrorIs(t, err, ErrInvalidConfig)
			assert.ErrorContains(t, err, "'"+key+"'")
			assert.Nil(t, result)
		})
	}

	t.Run("reads every v2.6.1 key", func(t *testing.T) {
		// given
		fileName := "../../testdata/config_full.json"
		// when
		result, err := Read(fileName)
		// then
		require.NoError(t, err)
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
	t.Run("a missing folder returns an error", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "missing", "config.json")
		// when
		err := Write(fileName, Default())
		// then
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("file reads back as the same configuration", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "config.json")
		// when
		err := Write(fileName, Default())
		// then
		require.NoError(t, err)
		result := mustRead(t, fileName)
		assert.Equal(t, Default(), result)
	})

	t.Run("creates the file readable only by its owner", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "config.json")
		// when
		err := Write(fileName, Default())
		// then
		require.NoError(t, err)
		info, err := os.Stat(fileName)
		require.NoError(t, err)
		expected := os.FileMode(0o600)
		assert.Equal(t, expected, info.Mode().Perm())
	})

	t.Run("file has exactly the v2.6.1 key names", func(t *testing.T) {
		// given
		fileName := filepath.Join(t.TempDir(), "config.json")
		// when
		err := Write(fileName, Default())
		// then
		require.NoError(t, err)
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

func mustRead(t testing.TB, fileName string) *Config {
	t.Helper()
	cfg, err := Read(fileName)
	if err != nil {
		t.Fatalf("Read(%s): unexpected error: %v", fileName, err)
	}
	return cfg
}
