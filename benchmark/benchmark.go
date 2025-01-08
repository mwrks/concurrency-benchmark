package benchmark

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

func Benchmark(url string, numRun int, numRequests int, concurrency int, ticketID int, initialStock int, fileName string) (int32, int32, int32, time.Duration, float64, time.Duration) {
	// Metrics
	var totalRequests int32
	var successfulRequests int32
	var failedRequests int32
	var totalLatency time.Duration

	// Worker pool
	var wg sync.WaitGroup
	requests := make(chan int, numRequests)
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requests {
				latency, err := MakeOrderRequest(url, ticketID)
				atomic.AddInt32(&totalRequests, 1)
				atomic.AddInt64((*int64)(&totalLatency), int64(latency))
				if err != nil {
					atomic.AddInt32(&failedRequests, 1)
					log.Printf("Request failed: %v\n", err)
				} else {
					atomic.AddInt32(&successfulRequests, 1)
				}
			}
		}()
	}

	// Rate control using a ticker (optional)
	// rate := 10000 // requests per second
	// ticker := time.NewTicker(time.Second / time.Duration(rate))
	// defer ticker.Stop()

	start := time.Now()

	// Enqueue requests
	go func() {
		for i := 0; i < numRequests; i++ {
			// <-ticker.C
			requests <- i
		}
		close(requests)
	}()

	// Wait for workers to finish
	wg.Wait()
	duration := time.Since(start)

	// Report results
	averageLatency := time.Duration(0)
	if totalRequests > 0 {
		averageLatency = totalLatency / time.Duration(totalRequests)
	}
	// Request per second
	rps := float64(totalRequests) / duration.Seconds()

	// Fetch orders
	orders, err := FetchOrders(ticketID)
	if err != nil {
		log.Fatalf("Failed to fetch orders: %v", err)
	}
	successfulOrders := len(orders)

	// Error rate
	mismatchRate := (float64(successfulOrders-initialStock) / float64(initialStock)) * 100

	// Log to Excel
	if err := LogToExcel(fileName, numRun, mismatchRate, initialStock, successfulOrders, orders, duration, rps, averageLatency); err != nil {
		log.Fatalf("Failed to log to Excel: %v", err)
	}
	log.Printf("Test results logged successfully to %s", fileName)
	return totalRequests, successfulRequests, failedRequests, duration, rps, averageLatency
}
