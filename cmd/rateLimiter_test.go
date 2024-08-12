package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	assert := assert.New(t)
	httpClient := &http.Client{}

	statusOK := 0
	statusTooManyReqs := 0
	for i := 0; i < 101; i++ {
		req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
		req.Header.Set("API_KEY", "abc123")

		response, err := httpClient.Do(req)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if response.StatusCode == http.StatusOK {
			statusOK++
		} else if response.StatusCode == http.StatusTooManyRequests {
			statusTooManyReqs++
		}
	}
	assert.Equal(100, statusOK)
	assert.Equal(1, statusTooManyReqs)

	statusOK = 0
	statusTooManyReqs = 0
	for i := 0; i < 11; i++ {
		req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

		response, err := httpClient.Do(req)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if response.StatusCode == http.StatusOK {
			statusOK++
		} else if response.StatusCode == http.StatusTooManyRequests {
			statusTooManyReqs++
		}
	}
	assert.Equal(10, statusOK)
	assert.Equal(1, statusTooManyReqs)
}
