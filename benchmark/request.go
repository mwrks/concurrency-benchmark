package benchmark

import (
	"bytes"
	"concurrency-benchmark/models"
	"concurrency-benchmark/utils"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func MakeOrderRequest(url string, ticketID int) (time.Duration, error) {
	order := models.Order{
		TicketID:  ticketID,
		OrderedBy: utils.GetRandomName(),
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
	resp, err := utils.HttpClient.Do(req)
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
