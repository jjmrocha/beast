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
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	// given
	expectedConfig := Default()
	// when
	cfg := Read("../testdata/empty.json")
	// then
	if !reflect.DeepEqual(cfg, expectedConfig) {
		t.Errorf("got %v expected %v", cfg, expectedConfig)
	}
}

func TestGetMaxIdleConnections(t *testing.T) {
	// given
	cfg := Default()
	parallelConns := 100
	var tests = []struct {
		maxIdleConnections int
		expected           int
	}{
		{0, parallelConns},
		{5, 5},
	}
	// then
	for _, test := range tests {
		cfg.MaxIdleConnections = test.maxIdleConnections
		result := cfg.GetMaxIdleConnections(parallelConns)
		if result != test.expected {
			t.Errorf("got %v expected %v", result, test.expected)
		}
	}
}
