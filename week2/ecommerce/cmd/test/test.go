package main

import (
	"bytes"
	"ecommerce/internal/order/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// OrderPayload represents the structure of the order to be sent in the POST request
type OrderPayload struct {
	UserID int                `json:"user_id"`
	Lines  []domain.OrderLine `json:"lines"`
}

// Create a custom HTTP client with increased timeouts and connection pooling
var client = &http.Client{
	Timeout: 30 * time.Second, // Increase timeout to 30 seconds
	Transport: &http.Transport{
		MaxIdleConns:        100,              // Max idle connections across all hosts
		MaxIdleConnsPerHost: 10,               // Max idle connections per host
		IdleConnTimeout:     30 * time.Second, // Idle connection timeout
	},
}

// Function to send POST request to the order API
func placeOrder(order OrderPayload, wg *sync.WaitGroup, resultChan chan string) {
	defer wg.Done()

	// Convert order payload to JSON
	jsonData, err := json.Marshal(order)
	if err != nil {
		resultChan <- fmt.Sprintf("Failed to marshal JSON: %v", err)
		return
	}

	// Send POST request using the custom HTTP client
	req, err := http.NewRequest("POST", "http://localhost:8080/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		resultChan <- fmt.Sprintf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		resultChan <- fmt.Sprintf("Request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		resultChan <- "Success: Order placed"
	} else {
		resultChan <- fmt.Sprintf("Failed: Status %d", resp.StatusCode)
	}
}

func main() {
	// Variables for concurrency control
	var wg sync.WaitGroup
	resultChan := make(chan string, 1000)

	// Start time for performance measurement
	start := time.Now()

	// Simulate 1000 concurrent requests
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		order := OrderPayload{
			UserID: 1,
			Lines: []domain.OrderLine{
				{
					ProductID: 1,
					Qty:       1,
				},
			},
		}

		// Launch a goroutine for each order
		go placeOrder(order, &wg, resultChan)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(resultChan)

	// Calculate elapsed time
	elapsed := time.Since(start)

	// Collect and display the results
	successCount := 0
	failCount := 0
	for result := range resultChan {
		fmt.Println(result)
		if result == "Success: Order placed" {
			successCount++
		} else {
			failCount++
		}
	}

	// Output summary
	fmt.Printf("\n--- Test Summary - Without Lock ---")
	fmt.Printf("\nTotal Requests: %d\n", successCount+failCount)
	fmt.Printf("\nTotal Success: %d\nTotal Failures: %d\n", successCount, failCount)
	fmt.Printf("Success Rate: %.2f%%\n", float64(successCount)/10)
	fmt.Printf("Test completed in %s\n", elapsed)
}
