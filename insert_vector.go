package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ FINNHUB_API_KEY not set.")
	}

	client := InitFinnhubClient()
	tickers := []string{"AAPL", "MSFT", "GOOG", "TSLA"}

	intervalPerRound := time.Duration(60/len(tickers)) * time.Second
	delayBetweenRequests := 15 * time.Millisecond

	for {
		start := time.Now()

		for _, ticker := range tickers {
			go func(t string) {
				quote, _, err := client.Quote(context.Background()).Symbol(t).Execute()
				if err != nil || quote.O == nil || quote.C == nil {
					log.Printf("⚠️ Failed to fetch quote for %s: %v", t, err)
					return
				}

				now := time.Now().UTC()
				tickerID := symbolToFloat(t)
				timeFloat := float64(now.Unix()) + float64(now.Nanosecond())/1e9

				vec := []float64{
					tickerID,
					timeFloat,
					float64(*quote.O),
					float64(*quote.H),
					float64(*quote.L),
					float64(*quote.C),
				}

				req := insertReq{
					Values:   vec,
					Metadata: map[string]string{"ticker": t, "time": now.Format(time.RFC3339)},
				}

				data, _ := json.Marshal(req)
				resp, err := http.Post("http://localhost:3000/insert", "application/json", bytes.NewBuffer(data))
				if err == nil {
					resp.Body.Close()
					log.Printf("✅ Inserted quote for %s at %s", t, now.Format("15:04:05.000"))
				}
			}(ticker)

			time.Sleep(delayBetweenRequests)
		}

		elapsed := time.Since(start)
		if elapsed < intervalPerRound {
			time.Sleep(intervalPerRound - elapsed)
		}
	}
}
