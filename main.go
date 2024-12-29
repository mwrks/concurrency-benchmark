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
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var paragraph = "Lorem ipsum dolor sit amet consectetur adipiscing elit Nullam vehicula ex id quam tincidunt ac varius justo cursus Proin ac efficitur risus quis dapibus tortor Praesent sit amet vehicula lorem vel pharetra mauris Aenean congue felis a sapien ultricies hendrerit Curabitur in sem vitae mi sagittis bibendum in nec elit Cras vel nisl vel risus dictum tincidunt vel id libero Sed aliquet dolor eget libero aliquet vel aliquet sem consequat Vivamus auctor justo in urna gravida faucibus Fusce luctus purus vel pharetra efficitur velit sapien tincidunt sapien eget vulputate quam ligula sed turpis Sed viverra hendrerit purus id posuere Ut quis finibus magna Aliquam sodales odio sed consequat maximus justo justo egestas lectus non commodo nisi sapien non ipsum Nullam non magna ut ligula accumsan fermentum Integer pellentesque velit eu orci aliquet id pharetra erat mollis Ut volutpat ligula nec ipsum fermentum sed interdum metus vehicula Suspendisse ac sapien at justo pharetra auctor in sed nisi Morbi molestie eros vel mauris tempor sodales Maecenas scelerisque erat id sapien aliquet vehicula Vestibulum scelerisque nisi sed rutrum scelerisque nisl enim aliquet dolor vel tincidunt sapien nulla et arcu Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas Ut at lectus non magna convallis blandit eget id tortor Curabitur gravida justo at lacinia dictum Duis malesuada lacinia quam nec cursus neque facilisis nec Donec efficitur suscipit tellus Quisque scelerisque orci et arcu vestibulum fermentum Fusce eget nulla nisl Cras vehicula sagittis tellus sit amet eleifend Praesent tincidunt sem ac tortor finibus quis mollis purus tincidunt Aenean tincidunt nunc vel tincidunt venenatis Sed vitae lectus id dolor dictum vehicula id nec sapien Nulla ac nunc nec enim interdum dictum in a libero Suspendisse potenti Praesent eget lacus nec sapien malesuada gravida in ac justo Duis fringilla justo et augue venenatis luctus Pellentesque consectetur ipsum quis velit bibendum non posuere nisi ultricies"
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
var names = strings.Fields(paragraph)

func getRandomName() string {
	return names[rnd.Intn(len(names))]
}

type Order struct {
	TicketID  int    `json:"ticket_id"`
	OrderedBy string `json:"ordered_by"`
}

// Shared HTTP client with connection pooling
var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

func makeOrderRequest(url string, ticketID int) error {
	order := Order{
		TicketID:  ticketID,
		OrderedBy: getRandomName(),
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Command-line flags for configuration
	url := flag.String("url", "http://localhost:3000/order", "Target URL")
	numRequests := flag.Int("requests", 10000, "Total number of requests to send")
	concurrency := flag.Int("concurrency", 100, "Number of concurrent workers")
	ticketID := flag.Int("ticket_id", 1, "Ticket ID to use for orders")
	flag.Parse()

	// Metrics
	var totalRequests int32
	var successfulRequests int32
	var failedRequests int32

	// Worker pool
	var wg sync.WaitGroup
	requests := make(chan int, *numRequests)
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requests {
				err := makeOrderRequest(*url, *ticketID)
				atomic.AddInt32(&totalRequests, 1)
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
	rate := 1000 // requests per second
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()

	start := time.Now()

	// Enqueue requests
	go func() {
		for i := 0; i < *numRequests; i++ {
			<-ticker.C
			requests <- i
		}
		close(requests)
	}()

	// Wait for workers to finish
	wg.Wait()
	duration := time.Since(start)

	// Report results
	fmt.Printf("Load Test Completed\n")
	fmt.Printf("===================\n")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Successful Requests: %d\n", successfulRequests)
	fmt.Printf("Failed Requests: %d\n", failedRequests)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Requests per Second: %.2f\n", float64(totalRequests)/duration.Seconds())
}
