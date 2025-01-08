package main

import (
	"concurrency-benchmark/benchmark"
	"concurrency-benchmark/utils"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	// Command-line flags for configuration
	fileName := flag.String("filename", "test_results.xlsx", "Test results filename")
	url := flag.String("url", "http://localhost:3000/order", "Target URL")
	numRequests := flag.Int("requests", 100, "Total number of requests to send")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent workers")
	ticketID := flag.Int("ticket_id", 1, "Ticket ID to use for orders")
	numRuns := flag.Int("runs", 1, "Number of times to run the tool")
	flag.Parse()

	// Check filename
	fileName = utils.FilenameCheck(fileName)

	fmt.Printf("Running Benchmark\n")
	fmt.Printf("Requests\t: %d\n", *numRequests)
	fmt.Printf("Concurrency\t: %d\n", *concurrency)
	fmt.Printf("Test runs\t: %d\n", *numRuns)
	fmt.Printf("Log file\t: %s\n", *fileName)
	fmt.Printf("===================\n")

	// Loop to run the test multiple times
	for i := 0; i < *numRuns; i++ {
		log.Printf("Running test %d of %d...", i+1, *numRuns)

		// Step 1: Get the initial stock
		initialStock, err := benchmark.FetchInitialStock(*ticketID)
		if err != nil {
			log.Fatalf("Failed to get initial stock: %v", err)
		}
		log.Printf("Initial stock for ticket ID %d: %d", *ticketID, initialStock)

		// Step 2: Update the current stock to match the initial stock
		if err := benchmark.UpdateCurrentStock(*ticketID, initialStock); err != nil {
			log.Fatalf("Failed to update current stock: %v", err)
		}
		log.Printf("Updated current stock for ticket ID %d to %d", *ticketID, initialStock)

		totalRequests, successfulRequests, failedRequests, duration, rps, averageLatency := benchmark.Benchmark(*url, i+1, *numRequests, *concurrency, *ticketID, initialStock, *fileName)

		// Step 3: Reset orders and sequence
		if err := benchmark.ResetOrders(*ticketID); err != nil {
			log.Printf("Failed to reset orders: %v", err)
		}

		// Checking orders is 0
		if orders, err := benchmark.FetchOrders(*ticketID); len(orders) != 0 {
			log.Fatalf("Orders table is not empty after reset: %v", err)
		}

		// Reset table serial sequence
		if err := benchmark.ResetOrderSequence(); err != nil {
			log.Printf("Failed to reset order sequence: %v", err)
		}

		fmt.Printf("\t- Total Requests\t: %d\n", totalRequests)
		fmt.Printf("\t- Successful Requests\t: %d\n", successfulRequests)
		fmt.Printf("\t- Failed Requests\t: %d\n", failedRequests)
		fmt.Printf("\t- Duration\t\t: %v\n", duration)
		fmt.Printf("\t- Requests per Second\t: %.2f\n", rps)
		fmt.Printf("\t- Average Latency\t: %v\n", averageLatency)
		// Add 2-second delay before the next run
		if i < *numRuns-1 {
			log.Printf("Waiting for 1 seconds before the next run...\n\n")
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("===================\n")
	fmt.Printf("Load test completed\n")
}
