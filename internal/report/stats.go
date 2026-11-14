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
		s.successMap[response.StatusCode] = append(s.successMap[response.StatusCode], response.Duration)
	case response.Failed():
		s.errorMap[response.FailureReason()]++
	default:
		s.statusMap[response.StatusCode]++
	}

	s.progress.Update()
}

// PrintStats displays the stats
func (s *Stats) PrintStats(w io.Writer) {
	s.printTotals(w)
	s.printSuccessStatuses(w)
	s.printNonSuccessStatuses(w)
	s.printErrors(w)
}

func (s *Stats) printTotals(w io.Writer) {
	_, _ = fmt.Fprintf(w, "===== Stats =====\n")
	_, _ = fmt.Fprintf(w, "Executed requests: %v\n", s.requests)
	_, _ = fmt.Fprintf(w, "Time taken to complete: %v\n", s.executionDuration())
	_, _ = fmt.Fprintf(w, "Requests per second: %.4f\n", s.tps())
	_, _ = fmt.Fprintf(w, "Avg response time: %v\n", avg(s.duration, s.requests))
}

func (s *Stats) printSuccessStatuses(w io.Writer) {
	for status, durations := range s.successMap {
		count := durations.Len()
		_, _ = fmt.Fprintf(w, "===== Status %v =====\n", status)
		_, _ = fmt.Fprintf(w, "%v requests, with avg response time of %v\n", count, avg(durations.sum(), count))

		if count >= minResultsForDistribution {
			printDistribution(w, durations)
		}
	}
}

const minResultsForDistribution = 5

var percentiles = []struct {
	percent    int
	minResults int
}{
	{percent: 20, minResults: minResultsForDistribution},
	{percent: 40, minResults: minResultsForDistribution},
	{percent: 60, minResults: minResultsForDistribution},
	{percent: 80, minResults: minResultsForDistribution},
	{percent: 90, minResults: 10},
	{percent: 95, minResults: 20},
	{percent: 99, minResults: 100},
}

func printDistribution(w io.Writer, durations durationSlice) {
	sort.Sort(durations)

	_, _ = fmt.Fprintf(w, "And the following distribution:\n")
	_, _ = fmt.Fprintf(w, "- The fastest request took %v\n", durations.first())
	for _, p := range percentiles {
		if durations.Len() >= p.minResults {
			_, _ = fmt.Fprintf(w, "- %v%% of requests under %v\n", p.percent, durations.percentage(p.percent))
		}
	}
	_, _ = fmt.Fprintf(w, "- The slowest request took %v\n", durations.last())
}

func (s *Stats) printNonSuccessStatuses(w io.Writer) {
	if len(s.statusMap) == 0 {
		return
	}

	_, _ = fmt.Fprintf(w, "===== Non Success Status =====\n")
	for status, count := range s.statusMap {
		_, _ = fmt.Fprintf(w, "Status %v: %v requests\n", status, count)
	}
}

func (s *Stats) printErrors(w io.Writer) {
	if len(s.errorMap) == 0 {
		return
	}

	_, _ = fmt.Fprintf(w, "===== Errors =====\n")
	for reason, count := range s.errorMap {
		_, _ = fmt.Fprintf(w, "- %v: %v errors\n", reason, count)
	}
}

func (s *Stats) tps() float64 {
	return float64(s.requests) / s.executionDuration().Seconds()
}

func (s *Stats) executionDuration() time.Duration {
	return time.Since(s.executionStart)
}

func avg(duration time.Duration, requests int) time.Duration {
	if requests == 0 {
		return 0
	}

	return time.Duration(duration.Nanoseconds() / int64(requests))
}
