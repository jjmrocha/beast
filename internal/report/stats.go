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
	"io"
	"sort"
	"time"

	"github.com/jjmrocha/beast/internal/client"
)

// Progress defines the progress indicator interface used by stats collector to inform user of the execution progress
type Progress interface {
	Update()
}

// Stats collects statistics about the results of the execution
type Stats struct {
	requests       int
	executionStart time.Time
	duration       time.Duration
	successMap     map[int]durationSlice
	statusMap      map[int]int
	errorMap       map[string]int
	progress       Progress
}

// NewStats creates a new Stats
func NewStats(progress Progress) *Stats {
	return &Stats{
		executionStart: time.Now(),
		successMap:     make(map[int]durationSlice),
		statusMap:      make(map[int]int),
		errorMap:       make(map[string]int),
		progress:       progress,
	}
}

// Update receives results and update the stats accordingly
func (s *Stats) Update(response *client.Response) {
	s.requests++
	s.duration += response.Duration

	switch {
	case response.IsSuccess():
		durations, present := s.successMap[response.StatusCode]
		if !present {
			durations = make(durationSlice, 0)
		}

		s.successMap[response.StatusCode] = append(durations, response.Duration)
	case response.IsClientError():
		errorDesc := response.ClientError()
		s.errorMap[errorDesc]++
	default:
		s.statusMap[response.StatusCode]++
	}

	s.progress.Update()
}

func (s *Stats) tps() float64 {
	return float64(s.requests) / s.executionDuration().Seconds()
}

func (s *Stats) avg() time.Duration {
	return avg(s.duration, s.requests)
}

func (s *Stats) executionDuration() time.Duration {
	return time.Since(s.executionStart)
}

func avg(duration time.Duration, requests int) time.Duration {
	return time.Duration(duration.Nanoseconds() / int64(requests))
}

// PrintStats displays the stats
func (s *Stats) PrintStats(w io.Writer) {
	_, _ = fmt.Fprintf(w, "===== Stats =====\n")
	_, _ = fmt.Fprintf(w, "Executed requests: %v\n", s.requests)
	_, _ = fmt.Fprintf(w, "Time taken to complete: %v\n", s.executionDuration())
	_, _ = fmt.Fprintf(w, "Requests per second: %.4f\n", s.tps())
	_, _ = fmt.Fprintf(w, "Avg response time: %v\n", s.avg())

	for key, durations := range s.successMap {
		_, _ = fmt.Fprintf(w, "===== Status %v =====\n", key)
		count := durations.Len()
		duration := durations.sum()
		_, _ = fmt.Fprintf(w, "%v requests, with avg response time of %v\n", count, avg(duration, count))
		if count >= 5 {
			sort.Sort(durations)
			_, _ = fmt.Fprintf(w, "And the following distribution:\n")
			_, _ = fmt.Fprintf(w, "- The fastest request took %v\n", durations.first())
			_, _ = fmt.Fprintf(w, "- 20%% of requests under %v\n", durations.percentage(20))
			_, _ = fmt.Fprintf(w, "- 40%% of requests under %v\n", durations.percentage(40))
			_, _ = fmt.Fprintf(w, "- 60%% of requests under %v\n", durations.percentage(60))
			_, _ = fmt.Fprintf(w, "- 80%% of requests under %v\n", durations.percentage(80))
			if count >= 10 {
				_, _ = fmt.Fprintf(w, "- 90%% of requests under %v\n", durations.percentage(90))
				if count >= 20 {
					_, _ = fmt.Fprintf(w, "- 95%% of requests under %v\n", durations.percentage(95))
					if count >= 100 {
						_, _ = fmt.Fprintf(w, "- 99%% of requests under %v\n", durations.percentage(99))
					}
				}
			}
			_, _ = fmt.Fprintf(w, "- The slowest request took %v\n", durations.last())
		}
	}

	if len(s.statusMap) > 0 {
		_, _ = fmt.Fprintf(w, "===== Non Success Status =====\n")

		for key, value := range s.statusMap {
			_, _ = fmt.Fprintf(w, "Status %v: %v requests\n", key, value)
		}
	}

	if len(s.errorMap) > 0 {
		_, _ = fmt.Fprintf(w, "===== Errors =====\n")

		for key, value := range s.errorMap {
			_, _ = fmt.Fprintf(w, "- %v: %v errors\n", key, value)
		}
	}

}
