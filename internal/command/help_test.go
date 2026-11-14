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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedTitle = "the Beast dev - Stress testing for RESTful APIs\n"

func TestVersion(t *testing.T) {
	// given
	var out strings.Builder
	// when
	Version(&out)
	// then
	assert.Equal(t, expectedTitle, out.String())
}

func TestHelp(t *testing.T) {
	// given
	var out strings.Builder
	// when
	Help(&out)
	// then
	result := out.String()
	expectedStart := expectedTitle + "\nUsage:\n"
	assert.True(t, strings.HasPrefix(result, expectedStart), "help output: %q", result)
	assert.Contains(t, result, "   beast version\n")
}
