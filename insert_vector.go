package main

import (
	"log"
	"os"
	"time"
)

func main() {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ FINNHUB_API_KEY not set.")
	}

	client := InitFinnhubClient()
	tickers := []string{"AAPL", "MSFT", "GOOG", "TSLA"} // Example: 4 tickers

	delayBetweenBatches := time.Duration(60/len(tickers)) * time.Second
	delayBetweenRequests := 15 * time.Millisecond

	for {
		start := time.Now()
		for _, ticker := range tickers {
			fetchAndInsertQuote(client, ticker)
			time.Sleep(delayBetweenRequests)
		}

		// wait remaining time until the next batch
		elapsed := time.Since(start)
		if elapsed < delayBetweenBatches {
			time.Sleep(delayBetweenBatches - elapsed)
		}
	}
}
