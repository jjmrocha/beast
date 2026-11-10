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

package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	t.Run("creates the example template when no URL is given", func(t *testing.T) {
		// given
		opts := TemplateOptions{File: filepath.Join(t.TempDir(), "request.yaml")}
		// when
		err := Template(opts)
		// then
		assert.NoError(t, err)
		assert.FileExists(t, opts.File)
	})

	t.Run("creates the template for the given URL", func(t *testing.T) {
		// given
		opts := TemplateOptions{Method: "POST", URL: "http://someendpoint.pt", File: filepath.Join(t.TempDir(), "request.json")}
		// when
		err := Template(opts)
		// then
		assert.NoError(t, err)
		assert.FileExists(t, opts.File)
	})

	t.Run("returns an error and creates no file when the folder does not exist", func(t *testing.T) {
		// given
		opts := TemplateOptions{File: filepath.Join(t.TempDir(), "missing", "request.yaml")}
		// when
		err := Template(opts)
		// then
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.NoFileExists(t, opts.File)
	})
}
