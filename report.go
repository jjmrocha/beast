package beast

import (
	"fmt"
	"time"
)

type Report struct {
	Requests  int
	Duration  time.Duration
	Min       time.Duration
	Max       time.Duration
	StatusMap map[int]int
}

func (o *Report) Update(r *BResponse) {
	response := *r
	o.Requests++
	o.Duration += response.Duration

	if o.Requests == 1 {
		o.Min = response.Duration
		o.Max = response.Duration
		o.StatusMap = make(map[int]int)
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
	return float64(o.Requests) / o.Duration.Seconds()
}

func (o *Report) Avg() time.Duration {
	return o.Duration / time.Duration(o.Requests)
}

func (o *Report) Print() {
	fmt.Printf("=== Results ===\n")
	fmt.Printf("Total requests: %v\n", o.Requests)
	fmt.Printf("Time taken to complete requests: %v\n", o.Duration)
	fmt.Printf("Requests per second: %.4f\n", o.Tps())
	fmt.Printf("Max response time: %v\n", o.Max)
	fmt.Printf("Min response time: %v\n", o.Min)
	fmt.Printf("Avg response time: %v\n", o.Avg())
	fmt.Printf("=== Status Count ===\n")

	for key, value := range o.StatusMap {
		fmt.Printf("Status %v: %v requests\n", key, value)
	}
}
