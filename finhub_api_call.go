package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type insertReq struct {
	Values   []float64         `json:"values"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type insertResp struct {
	UUID string `json:"uuid"`
}

func main() {
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ FINNHUB_API_KEY environment variable not set.")
	}

	client := InitFinnhubClient()
	ticker := "AAPL"

	vector, tradingDay, err := FetchQuoteVector(client, ticker)
	if err != nil {
		log.Fatalf("❌ Error: %v", err)
	}

	// Convert vector from float32[] to float64[] for JSON
	values := make([]float64, len(vector))
	for i, v := range vector {
		values[i] = float64(v)
	}

	reqBody := insertReq{
		Values:   values,
		Metadata: map[string]string{"ticker": ticker, "date": tradingDay.Format("2006-01-02")},
	}
	payload, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:3000/insert", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("❌ Failed to send insert request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("❌ Insert failed with status %d", resp.StatusCode)
	}

	var insertResponse insertResp
	if err := json.NewDecoder(resp.Body).Decode(&insertResponse); err != nil {
		log.Fatalf("❌ Failed to parse response: %v", err)
	}

	fmt.Printf("✅ Inserted vector for %s on %s → UUID: %s\n", ticker, tradingDay.Format("2006-01-02"), insertResponse.UUID)
}
