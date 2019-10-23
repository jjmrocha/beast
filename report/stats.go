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

package report

import (
	"fmt"
	"sort"
	"time"

	"github.com/jjmrocha/beast/client"
)

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

// Stats collects statistics about the results of the execution
type Stats struct {
	concurrent int
	requests   int
	duration   time.Duration
	successMap map[int]durationSlice
	errorMap   map[int]int
}

// NewStats creates a new Stats
func NewStats(nParallel int) *Stats {
	return &Stats{
		concurrent: nParallel,
		successMap: make(map[int]durationSlice),
		errorMap:   make(map[int]int),
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
	} else {
		s.errorMap[r.StatusCode]++
	}
}

func success(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
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
	fmt.Printf("=== Result Stats ===\n")
	fmt.Printf("Executed requests: %v\n", s.requests)
	fmt.Printf("Time taken to complete: %v\n", s.duration)
	fmt.Printf("Requests per second: %.4f\n", s.tps())
	fmt.Printf("Avg response time: %v\n", s.avg())

	for key, durations := range s.successMap {
		fmt.Printf("=== Status %v ===\n", key)
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

	if len(s.errorMap) > 0 {
		fmt.Printf("=== Non Success Status ===\n")

		for key, value := range s.errorMap {
			fmt.Printf("Status %v: %v requests\n", key, value)
		}
	}
}
