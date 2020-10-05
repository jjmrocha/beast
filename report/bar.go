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
	"fmt"
	"log"
	"strings"
	"time"
)

// Bar represents a progress bar
type Bar struct {
	requestCount            int
	executionDuration       int
	executionStart          time.Time
	executedRequestsCount   int
	lastPercentageDisplayed int
}

// NewBar creates a new progress bar
func NewBar(max, duration int) *Bar {
	return &Bar{
		requestCount:      max,
		executionDuration: duration,
		executionStart:    time.Now(),
	}
}

// Update indicates to the progress bar that we receive another output,
// the bar will update its representation accordingly
func (b *Bar) Update() {
	b.executedRequestsCount++
	percentage := b.percentage()

	if percentage > 100 {
		percentage = 100
	}

	if percentage != b.lastPercentageDisplayed && percentage%5 == 0 {
		b.lastPercentageDisplayed = percentage
		bar := drawBar(percentage)
		log.Printf("%s %v%%\n", bar, percentage)
	}
}

func (b *Bar) percentage() int {
	if b.requestCount > 0 {
		return b.executedRequestsCount * 100.0 / b.requestCount
	}

	return int(time.Since(b.executionStart).Seconds()) * 100.0 / b.executionDuration
}

func drawBar(percentage int) string {
	count := percentage / 5
	done := strings.Repeat("#", count)
	todo := strings.Repeat(".", 20-count)
	return fmt.Sprintf("[%s%s]", done, todo)
}
