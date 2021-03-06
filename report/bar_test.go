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

import "testing"

func TestDrawBar(t *testing.T) {
	// given
	var tests = []struct {
		input    int
		expected string
	}{
		{0, "[....................]"},
		{5, "[#...................]"},
		{10, "[##..................]"},
		{15, "[###.................]"},
		{20, "[####................]"},
		{25, "[#####...............]"},
		{30, "[######..............]"},
		{35, "[#######.............]"},
		{40, "[########............]"},
		{45, "[#########...........]"},
		{50, "[##########..........]"},
		{55, "[###########.........]"},
		{60, "[############........]"},
		{65, "[#############.......]"},
		{70, "[##############......]"},
		{75, "[###############.....]"},
		{80, "[################....]"},
		{85, "[#################...]"},
		{90, "[##################..]"},
		{95, "[###################.]"},
		{100, "[####################]"},
	}
	// then
	for _, test := range tests {
		result := drawBar(test.input)
		if result != test.expected {
			t.Errorf("got %v expected %v", result, test.expected)
		}
	}
}
