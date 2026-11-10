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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	const (
		success      = 0
		usageError   = 2
		runCommand   = "run"
		templateFile = "x.yaml"
	)

	tests := []struct {
		name     string
		input    []string
		expected int
	}{
		{name: "no arguments prints help", input: nil, expected: success},
		{name: "help", input: []string{"help"}, expected: success},
		{name: "unknown command", input: []string{"bogus"}, expected: usageError},
		{name: "config without a file", input: []string{"config"}, expected: usageError},
		{name: "run without -n or -t", input: []string{runCommand, templateFile}, expected: usageError},
		{name: "run with both -n and -t", input: []string{runCommand, "-n", "1", "-t", "1", templateFile}, expected: usageError},
		{name: "run with -c 0", input: []string{runCommand, "-n", "1", "-c", "0", templateFile}, expected: usageError},
		{name: "run without a template", input: []string{runCommand, "-n", "1"}, expected: usageError},
		{name: "template with three arguments", input: []string{"template", "a", "b", "c"}, expected: usageError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			args := tc.input
			// when
			result := run(args)
			// then
			assert.Equal(t, tc.expected, result)
		})
	}
}
