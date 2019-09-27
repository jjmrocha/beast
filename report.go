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
	"time"
)

type Report struct {
	concurrent int
	Requests   int
	Duration   time.Duration
	Min        time.Duration
	Max        time.Duration
	StatusMap  map[int]int
}

func NewReport(nParallel int) *Report {
	return &Report{
		concurrent: nParallel,
		StatusMap:  make(map[int]int),
	}
}

func (o *Report) Update(r *BResponse) {
	response := *r
	o.Requests++
	o.Duration += response.Duration

	if o.Requests == 1 {
		o.Min = response.Duration
		o.Max = response.Duration
	} else {
		if o.Min > response.Duration {
			o.Min = response.Duration
		}

		if o.Max < response.Duration {
			o.Max = response.Duration
		}
	}

	o.StatusMap[response.StatusCode]++
}

func (o *Report) Tps() float64 {
	return float64(o.concurrent) * (float64(o.Requests) / o.Duration.Seconds())
}

func (o *Report) Avg() time.Duration {
	return o.Duration / time.Duration(o.Requests)
}

func (o *Report) Print() {
	fmt.Printf("=== Results ===\n")
	fmt.Printf("Executed requests: %v\n", o.Requests)
	fmt.Printf("Time taken to complete: %v\n", o.Duration)
	fmt.Printf("=== Stats ===\n")
	fmt.Printf("Min response time: %v\n", o.Min)
	fmt.Printf("Max response time: %v\n", o.Max)
	fmt.Printf("Avg response time: %v\n", o.Avg())
	fmt.Printf("Requests per second: %.4f\n", o.Tps())
	fmt.Printf("=== Status Code ===\n")

	for key, value := range o.StatusMap {
		fmt.Printf("Status %v: %v requests\n", key, value)
	}
}
