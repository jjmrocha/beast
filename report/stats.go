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

type Stats struct {
	concurrent int
	Requests   int
	Duration   time.Duration
	SuccessMap map[int]durationSlice
	ErrorMap   map[int]int
}

func NewStats(nParallel int) *Stats {
	return &Stats{
		concurrent: nParallel,
		SuccessMap: make(map[int]durationSlice),
		ErrorMap:   make(map[int]int),
	}
}

func (s *Stats) Update(r *client.BResponse) {
	s.Requests++
	s.Duration += r.Duration

	if success(r.StatusCode) {
		durations, present := s.SuccessMap[r.StatusCode]
		if !present {
			durations = make(durationSlice, 0)
		}

		s.SuccessMap[r.StatusCode] = append(durations, r.Duration)
	} else {
		s.ErrorMap[r.StatusCode]++
	}
}

func success(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func (s *Stats) Tps() float64 {
	return float64(s.concurrent) * (float64(s.Requests) / s.Duration.Seconds())
}

func (s *Stats) Avg() time.Duration {
	return avg(s.Duration, s.Requests)
}

func avg(duration time.Duration, requests int) time.Duration {
	return time.Duration(duration.Nanoseconds() / int64(requests))
}

func (s *Stats) Print() {
	fmt.Printf("=== Results ===\n")
	fmt.Printf("Executed requests: %v\n", s.Requests)
	fmt.Printf("Time taken to complete: %v\n", s.Duration)
	fmt.Printf("Requests per second: %.4f\n", s.Tps())
	fmt.Printf("Avg response time: %v\n", s.Avg())

	for key, durations := range s.SuccessMap {
		fmt.Printf("=== Status %v ===\n", key)
		count := len(durations)
		duration := sum(durations)
		sort.Sort(durations)
		fmt.Printf("%v requests, with avg response time of %v\n", count, avg(duration, count))
		if count >= 5 {
			fmt.Printf("And the following distribution:\n")
			fmt.Printf("  The fastest request took %v\n", durations[0])
			fmt.Printf("  20%% of requests under %v\n", durations[count/5-1])
			fmt.Printf("  40%% of requests under %v\n", durations[count/5*2-1])
			fmt.Printf("  60%% of requests under %v\n", durations[count/5*3-1])
			fmt.Printf("  80%% of requests under %v\n", durations[count/5*4-1])
			fmt.Printf("  The slowest request took %v\n", durations[count-1])
		}
	}

	if len(s.ErrorMap) > 0 {
		fmt.Printf("=== Errors ===\n")

		for key, value := range s.ErrorMap {
			fmt.Printf("Status %v: %v requests\n", key, value)
		}
	}
}

func sum(durations durationSlice) time.Duration {
	var sum time.Duration

	for _, value := range durations {
		sum += value
	}

	return sum
}
