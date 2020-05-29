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
	"sort"
	"time"

	"github.com/jjmrocha/beast/client"
)

// Stats collects statistics about the results of the execution
type Stats struct {
	concurrent int
	requests   int
	duration   time.Duration
	successMap map[int]durationSlice
	statusMap  map[int]int
	errorMap   map[errorCode]int
}

// NewStats creates a new Stats
func NewStats(nParallel int) *Stats {
	return &Stats{
		concurrent: nParallel,
		successMap: make(map[int]durationSlice),
		statusMap:  make(map[int]int),
		errorMap:   make(map[errorCode]int),
	}
}

// Update receives results and update the stats accordingly
func (s *Stats) Update(r *client.BResponse) {
	s.requests++
	s.duration += r.Duration

	if success(r.StatusCode) {
		durations, present := s.successMap[r.StatusCode]
		if !present {
			durations = make(durationSlice, 0)
		}

		s.successMap[r.StatusCode] = append(durations, r.Duration)
	} else if error(r.StatusCode) {
		errorCode := errorCode(r.StatusCode)
		s.errorMap[errorCode]++
	} else {
		s.statusMap[r.StatusCode]++
	}
}

func success(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func error(statusCode int) bool {
	return statusCode < 0
}

func (s *Stats) tps() float64 {
	return float64(s.concurrent) * (float64(s.requests) / s.duration.Seconds())
}

func (s *Stats) avg() time.Duration {
	return avg(s.duration, s.requests)
}

func avg(duration time.Duration, requests int) time.Duration {
	return time.Duration(duration.Nanoseconds() / int64(requests))
}

// Print displays the stats
func (s *Stats) Print() {
	fmt.Printf("===== Stats =====\n")
	fmt.Printf("Executed requests: %v\n", s.requests)
	fmt.Printf("Time taken to complete: %v\n", s.duration)
	fmt.Printf("Requests per second: %.4f\n", s.tps())
	fmt.Printf("Avg response time: %v\n", s.avg())

	for key, durations := range s.successMap {
		fmt.Printf("===== Status %v =====\n", key)
		count := durations.Len()
		duration := durations.sum()
		fmt.Printf("%v requests, with avg response time of %v\n", count, avg(duration, count))
		if count >= 5 {
			sort.Sort(durations)
			fmt.Printf("And the following distribution:\n")
			fmt.Printf("- The fastest request took %v\n", durations.first())
			fmt.Printf("- 20%% of requests under %v\n", durations.percentage(20))
			fmt.Printf("- 40%% of requests under %v\n", durations.percentage(40))
			fmt.Printf("- 60%% of requests under %v\n", durations.percentage(60))
			fmt.Printf("- 80%% of requests under %v\n", durations.percentage(80))
			fmt.Printf("- The slowest request took %v\n", durations.last())
		}
	}

	if len(s.statusMap) > 0 {
		fmt.Printf("===== Non Success Status =====\n")

		for key, value := range s.statusMap {
			fmt.Printf("Status %v: %v requests\n", key, value)
		}
	}

	if len(s.errorMap) > 0 {
		fmt.Printf("===== Errors =====\n")

		for key, value := range s.errorMap {
			fmt.Printf("- %v: %v errors\n", key, value)
		}
	}
}
