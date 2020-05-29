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

import "time"

type durationSlice []time.Duration

func (a durationSlice) Len() int {
	return len(a)
}

func (a durationSlice) Less(i, j int) bool {
	return a[i] < a[j]
}

func (a durationSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a durationSlice) first() time.Duration {
	return a[0]
}

func (a durationSlice) last() time.Duration {
	return a[len(a)-1]
}

func (a durationSlice) percentage(value int) time.Duration {
	size := len(a)
	pos := (size * value / 100) - 1
	return a[pos]
}

func (a durationSlice) sum() time.Duration {
	var sum time.Duration

	for _, value := range a {
		sum += value
	}

	return sum
}
