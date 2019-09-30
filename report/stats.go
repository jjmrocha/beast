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
	"time"

	"github.com/jjmrocha/beast/client"
)

type Stats struct {
	concurrent int
	Requests   int
	Duration   time.Duration
	Min        time.Duration
	Max        time.Duration
	StatusMap  map[int]int
}

func NewReport(nParallel int) *Stats {
	return &Stats{
		concurrent: nParallel,
		StatusMap:  make(map[int]int),
	}
}

func (s *Stats) Update(r *client.BResponse) {
	response := *r
	s.Requests++
	s.Duration += response.Duration

	if s.Requests == 1 {
		s.Min = response.Duration
		s.Max = response.Duration
	} else {
		if s.Min > response.Duration {
			s.Min = response.Duration
		}

		if s.Max < response.Duration {
			s.Max = response.Duration
		}
	}

	s.StatusMap[response.StatusCode]++
}

func (s *Stats) Tps() float64 {
	return float64(s.concurrent) * (float64(s.Requests) / s.Duration.Seconds())
}

func (s *Stats) Avg() time.Duration {
	return time.Duration(s.Duration.Nanoseconds() / int64(s.Requests))
}

func (s *Stats) Print() {
	fmt.Printf("=== Results ===\n")
	fmt.Printf("Executed requests: %v\n", s.Requests)
	fmt.Printf("Time taken to complete: %v\n", s.Duration)
	fmt.Printf("=== Stats ===\n")
	fmt.Printf("Min response time: %v\n", s.Min)
	fmt.Printf("Max response time: %v\n", s.Max)
	fmt.Printf("Avg response time: %v\n", s.Avg())
	fmt.Printf("Requests per second: %.4f\n", s.Tps())
	fmt.Printf("=== Status Code ===\n")

	for key, value := range s.StatusMap {
		fmt.Printf("Status %v: %v requests\n", key, value)
	}
}
