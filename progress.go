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
package beast

import (
	"fmt"
	"log"
	"strings"
)

type Bar struct {
	Max     int
	Current int
	last    int
}

func NewBar(max int) *Bar {
	return &Bar{
		Max: max,
	}
}

func (b *Bar) Update() {
	b.Current++
	percentage := b.Current * 100.0 / b.Max

	if percentage != b.last && percentage%5 == 0 {
		b.last = percentage
		bar := drawBar(percentage)
		log.Printf("%s %v%%\n", bar, percentage)
	}
}

func drawBar(percentage int) string {
	count := percentage / 5
	done := strings.Repeat("#", count)
	todo := strings.Repeat(".", 20-count)
	return fmt.Sprintf("[%s%s]", done, todo)
}
