package utils

import (
	"net/http"
	"time"
)

// Shared HTTP client with connection pooling
var HttpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}
