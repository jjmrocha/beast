package beast

import (
	"io"
	"net/http"
	"strings"
	"time"
)

func HttpClient() *http.Client {
	transport := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives:  false,
	}

	return &http.Client{Transport: transport}
}

func (r *BRequest) toReq() (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.Endpoint, bodyReader(r.Body))
	if err != nil {
		return nil, err
	}

	for _, header := range r.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	return req, nil
}

func bodyReader(body string) io.Reader {
	if body == "" {
		return nil
	}

	return strings.NewReader(body)
}

func Execute(client *http.Client, request *BRequest) *BResponse {
	req, err := request.toReq()
	if err != nil {
		return &BResponse{
			StatusCode: -1,
			Duration:   time.Duration(0),
		}
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return &BResponse{
			StatusCode: -2,
			Duration:   duration,
		}
	}

	return &BResponse{
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}
