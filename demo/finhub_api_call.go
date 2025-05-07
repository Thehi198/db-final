package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	apiKey := os.Getenv("FINNHUB_API_KEY") // Corrected env variable
	if apiKey == "" {
		log.Fatal("Please set your FINNHUB_API_KEY environment variable.") // Corrected error message
	}

	ticker := "AMD" // Example ticker

	vec, dateStr, err := FetchLatestVector(ticker, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Embedded vector for latest available date", dateStr, ":", vec)
}

// test commit
