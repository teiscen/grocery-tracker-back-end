//go:build ignore

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"receiptTracker/database"
)

// ScanService handles all receipt scanning logic
// MLBaseURL is the address of the ML VM's Crow endpoint
type ScanService struct {
	MLBaseURL 	string
	DB         *database.DB
}

// ReceiptItem represents a single line item parsed from a receipt
type ReceiptItem struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

// ReceiptResult is the full parsed receipt returned by the ML VM
type ReceiptResult struct {
	Vendor string        `json:"vendor"`
	Date   string        `json:"date"`
	Items  []ReceiptItem `json:"items"`
	Total  float64       `json:"total"`
}

// Process takes an image file, forwards it to the ML VM, and returns
// a parsed ReceiptResult. Returns an error if anything goes wrong.
func (s *ScanService) Process(file io.Reader) (*ReceiptResult, error) {
	// if no ML URL is configured, return mock data
	// this lets us develop without the ML VM running
	if s.MLBaseURL == "" {
		return &ReceiptResult{
			Vendor: "Superstore",
			Date:   "2026-05-11",
			Items: []ReceiptItem{
				{Name: "Whole Milk 1L", Quantity: 2, Price: 3.99},
				{Name: "Sourdough Bread", Quantity: 1, Price: 4.49},
				{Name: "Cheddar Cheese 400g", Quantity: 1, Price: 7.99},
			},
			Total: 20.46,
		}, nil
	}
	// forward the image to the ML VM
	result, err := s.callMLService(file)
	if err != nil {
		return nil, fmt.Errorf("ml service error: %w", err)
	}

	return result, nil
}

// callMLService forwards the image to the ML VM and parses the response
func (s *ScanService) callMLService(file io.Reader) (*ReceiptResult, error) {
	// build a multipart form to send to the ML VM
	// same format the client used to send to us
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("image", "receipt.jpg")
	if err != nil {
		return nil, err
	}

	// copy the image bytes into the form
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	writer.Close()

	// send the request to the ML VM
	resp, err := http.Post(
		s.MLBaseURL+"/ocr",
		writer.FormDataContentType(),
		&buf,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reach ml vm: %w", err)
	}
	defer resp.Body.Close()

	// check the ML VM returned a success status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml vm returned status %d", resp.StatusCode)
	}

	// parse the JSON response from the ML VM
	var result ReceiptResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse ml response: %w", err)
	}

	return &result, nil
}

