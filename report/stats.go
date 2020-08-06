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

// Progress defines the progress indicator interface used by stats collector to inform user of the execution progress
type Progress interface {
	Update()
}

// Stats collects statistics about the results of the execution
type Stats struct {
	concurrent     int
	requests       int
	executionStart time.Time
	duration       time.Duration
	successMap     map[int]durationSlice
	statusMap      map[int]int
	errorMap       map[string]int
	progress       Progress
	output         Output
}

// NewStats creates a new Stats
func NewStats(nParallel int, progress Progress, outputFile string) Stats {
	return Stats{
		concurrent:     nParallel,
		executionStart: time.Now(),
		successMap:     make(map[int]durationSlice),
		statusMap:      make(map[int]int),
		errorMap:       make(map[string]int),
		progress:       progress,
		output:         NewOutput(outputFile),
	}
}

// Update receives results and update the stats accordingly
func (s Stats) Update(r *client.Response) {
	s.requests++
	s.duration += r.Duration

	if r.IsSuccess() {
		durations, present := s.successMap[r.StatusCode]
		if !present {
			durations = make(durationSlice, 0)
		}

		s.successMap[r.StatusCode] = append(durations, r.Duration)
	} else if r.IsClientError() {
		errorDesc := r.ClientError()
		s.errorMap[errorDesc]++
	} else {
		s.statusMap[r.StatusCode]++
	}

	s.progress.Update()
	s.output.Write(r)
}

func (s Stats) tps() float64 {
	return float64(s.concurrent) * (float64(s.requests) / s.duration.Seconds())
}

func (s Stats) avg() time.Duration {
	return avg(s.duration, s.requests)
}

func (s Stats) executionDuration() time.Duration {
	return time.Since(s.executionStart)
}

func avg(duration time.Duration, requests int) time.Duration {
	return time.Duration(duration.Nanoseconds() / int64(requests))
}

// PrintStats displays the stats
func (s Stats) PrintStats() {
	fmt.Printf("===== Stats =====\n")
	fmt.Printf("Executed requests: %v\n", s.requests)
	fmt.Printf("Time taken to complete: %v\n", s.executionDuration())
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
			if count >= 10 {
				fmt.Printf("- 90%% of requests under %v\n", durations.percentage(90))
				if count >= 20 {
					fmt.Printf("- 95%% of requests under %v\n", durations.percentage(95))
					if count >= 100 {
						fmt.Printf("- 99%% of requests under %v\n", durations.percentage(99))
					}
				}
			}
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

	s.output.Close()
}
