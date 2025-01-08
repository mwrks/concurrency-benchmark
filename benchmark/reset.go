package benchmark

import (
	"bytes"
	"concurrency-benchmark/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UpdateCurrentStock(ticketID int, currentStock int) error {
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

	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func ResetOrders(ticketID int) error {
	url := fmt.Sprintf("http://localhost:3000/order/%d/reset", ticketID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

func ResetOrderSequence() error {
	url := "http://localhost:3000/order/reset-sequence"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}
