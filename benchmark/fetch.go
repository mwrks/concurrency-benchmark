package benchmark

import (
	"concurrency-benchmark/models"
	"concurrency-benchmark/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchInitialStock(ticketID int) (int, error) {
	url := fmt.Sprintf("http://localhost:3000/ticket/%d", ticketID)
	resp, err := utils.HttpClient.Get(url)
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

func FetchOrders(ticketID int) ([]models.Order, error) {
	url := fmt.Sprintf("http://localhost:3000/order/%d", ticketID)

	// Send the GET request to fetch orders
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the JSON response into the orders slice
	var orders []models.Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return orders, nil
}
