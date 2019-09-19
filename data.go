package beast

import "time"

type BHeader struct {
	Key   string
	Value string
}

type BRequest struct {
	Method   string
	Endpoint string
	Headers  []BHeader
	Body     string
}

type BResponse struct {
	StatusCode int
	Duration   time.Duration
}

type THeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TRequest struct {
	Method   string    `json:"method"`
	Endpoint string    `json:"url"`
	Headers  []THeader `json:"headers"`
	Body     string    `json:"body"`
}
