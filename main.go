package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xuri/excelize/v2"
)

var paragraph = "Lorem ipsum dolor sit amet consectetur adipiscing elit Nullam vehicula ex id quam tincidunt ac varius justo cursus Proin ac efficitur risus quis dapibus tortor Praesent sit amet vehicula lorem vel pharetra mauris Aenean congue felis a sapien ultricies hendrerit Curabitur in sem vitae mi sagittis bibendum in nec elit Cras vel nisl vel risus dictum tincidunt vel id libero Sed aliquet dolor eget libero aliquet vel aliquet sem consequat Vivamus auctor justo in urna gravida faucibus Fusce luctus purus vel pharetra efficitur velit sapien tincidunt sapien eget vulputate quam ligula sed turpis Sed viverra hendrerit purus id posuere Ut quis finibus magna Aliquam sodales odio sed consequat maximus justo justo egestas lectus non commodo nisi sapien non ipsum Nullam non magna ut ligula accumsan fermentum Integer pellentesque velit eu orci aliquet id pharetra erat mollis Ut volutpat ligula nec ipsum fermentum sed interdum metus vehicula Suspendisse ac sapien at justo pharetra auctor in sed nisi Morbi molestie eros vel mauris tempor sodales Maecenas scelerisque erat id sapien aliquet vehicula Vestibulum scelerisque nisi sed rutrum scelerisque nisl enim aliquet dolor vel tincidunt sapien nulla et arcu Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas Ut at lectus non magna convallis blandit eget id tortor Curabitur gravida justo at lacinia dictum Duis malesuada lacinia quam nec cursus neque facilisis nec Donec efficitur suscipit tellus Quisque scelerisque orci et arcu vestibulum fermentum Fusce eget nulla nisl Cras vehicula sagittis tellus sit amet eleifend Praesent tincidunt sem ac tortor finibus quis mollis purus tincidunt Aenean tincidunt nunc vel tincidunt venenatis Sed vitae lectus id dolor dictum vehicula id nec sapien Nulla ac nunc nec enim interdum dictum in a libero Suspendisse potenti Praesent eget lacus nec sapien malesuada gravida in ac justo Duis fringilla justo et augue venenatis luctus Pellentesque consectetur ipsum quis velit bibendum non posuere nisi ultricies"
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
var names = strings.Fields(paragraph)

func getRandomName() string {
	return names[rnd.Intn(len(names))]
}

type Order struct {
	OrderID   int       `json:"order_id"`
	TicketID  int       `json:"ticket_id"`
	OrderedBy string    `json:"ordered_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Shared HTTP client with connection pooling
var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

func getInitialStock(ticketID int) (int, error) {
	url := fmt.Sprintf("http://localhost:3000/ticket/%d", ticketID)
	resp, err := httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch initial stock: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		InitialStock int `json:"initial_quantity"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.InitialStock, nil
}

func updateCurrentStock(ticketID int, currentStock int) error {
	url := fmt.Sprintf("http://localhost:3000/ticket/%d", ticketID)

	body := map[string]int{"current_quantity": currentStock}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func makeOrderRequest(url string, ticketID int) (time.Duration, error) {
	order := Order{
		TicketID:  ticketID,
		OrderedBy: getRandomName(),
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := httpClient.Do(req)
	duration := time.Since(start)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

func resetOrders(ticketID int) error {
	url := fmt.Sprintf("http://localhost:3000/order/%d/reset", ticketID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

func resetOrderSequence() error {
	url := "http://localhost:3000/order/reset-sequence"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

func fetchOrders(ticketID int) ([]Order, error) {
	url := fmt.Sprintf("http://localhost:3000/order/%d", ticketID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var orders []Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return orders, nil
}

func logToExcel(fileName string, initialStock int, successfulOrders int, jsonBody []Order, duration time.Duration, rps float64, avgLatency time.Duration) error {
	// Create or open the Excel file
	var f *excelize.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// File doesn't exist, create a new one
		f = excelize.NewFile()
		// Create headers
		headers := []string{"no_run", "initial_stock", "amount_of_successful_orders", "json_body", "duration", "requests_per_second", "average_latency"}
		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue("Sheet1", cell, header)
		}
	} else {
		// Open existing file
		var err error
		f, err = excelize.OpenFile(fileName)
		if err != nil {
			return fmt.Errorf("failed to open Excel file: %v", err)
		}
	}

	// Find the next available row
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("failed to get rows: %v", err)
	}
	rowNo := len(rows) + 1

	// Convert JSON body to string
	jsonData, err := json.MarshalIndent(jsonBody, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Write data
	data := []interface{}{
		initialStock,
		successfulOrders,
		string(jsonData),
		duration.Seconds(),
		rps,
		avgLatency.Seconds(),
	}
	for i, value := range data {
		cell, _ := excelize.CoordinatesToCellName(i+1, rowNo)
		f.SetCellValue("Sheet1", cell, value)
	}

	// Save the file
	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("failed to save Excel file: %v", err)
	}
	return nil
}

func main() {
	// Command-line flags for configuration
	fileName := "test_results.xlsx"
	url := flag.String("url", "http://localhost:3000/order", "Target URL")
	numRequests := flag.Int("requests", 100, "Total number of requests to send")
	concurrency := flag.Int("concurrency", 100, "Number of concurrent workers")
	ticketID := flag.Int("ticket_id", 1, "Ticket ID to use for orders")
	numRuns := flag.Int("runs", 1, "Number of times to run the tool")
	flag.Parse()

	fmt.Printf("Running Benchmark\n")
	fmt.Printf("Requests\t: %d\n", *numRequests)
	fmt.Printf("Concurrency\t: %d\n", *concurrency)
	fmt.Printf("Test runs\t: %d\n", *numRuns)
	fmt.Printf("===================\n")
	// Loop to run the test multiple times
	for i := 0; i < *numRuns; i++ {
		log.Printf("Running test %d of %d...", i+1, *numRuns)

		// Step 1: Get the initial stock
		initialStock, err := getInitialStock(*ticketID)
		if err != nil {
			log.Fatalf("Failed to get initial stock: %v", err)
		}
		log.Printf("Initial stock for ticket ID %d: %d", *ticketID, initialStock)

		// Step 2: Update the current stock to match the initial stock
		if err := updateCurrentStock(*ticketID, initialStock); err != nil {
			log.Fatalf("Failed to update current stock: %v", err)
		}
		log.Printf("Updated current stock for ticket ID %d to %d", *ticketID, initialStock)

		// Metrics
		var totalRequests int32
		var successfulRequests int32
		var failedRequests int32
		var totalLatency time.Duration

		// Worker pool
		var wg sync.WaitGroup
		requests := make(chan int, *numRequests)
		for i := 0; i < *concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range requests {
					latency, err := makeOrderRequest(*url, *ticketID)
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
			for i := 0; i < *numRequests; i++ {
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
		orders, err := fetchOrders(*ticketID)
		if err != nil {
			log.Fatalf("Failed to fetch orders: %v", err)
		}
		successfulOrders := len(orders)

		// Log to Excel
		if err := logToExcel(fileName, initialStock, successfulOrders, orders, duration, rps, averageLatency); err != nil {
			log.Fatalf("Failed to log to Excel: %v", err)
		}

		log.Printf("Test results logged successfully to %s", fileName)
		// Step 3: Reset orders and sequence
		if err := resetOrders(*ticketID); err != nil {
			log.Printf("Failed to reset orders: %v", err)
		}

		if err := resetOrderSequence(); err != nil {
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
