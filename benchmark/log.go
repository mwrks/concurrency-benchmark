package benchmark

import (
	"concurrency-benchmark/models"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

func LogToExcel(fileName string, numRun int, errorRate int, initialStock int, successfulOrders int, jsonBody []models.Order, duration time.Duration, rps float64, avgLatency time.Duration) error {
	// Create or open the Excel file
	var f *excelize.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// File doesn't exist, create a new one
		f = excelize.NewFile()
		// Create headers
		headers := []string{"run", "error_rate", "initial_stock", "amount_of_successful_orders", "json_body", "duration", "requests_per_second", "average_latency"}
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
		numRun,
		errorRate,
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
