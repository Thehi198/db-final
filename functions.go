package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

type insertReq struct {
	Values   []float64         `json:"values"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func InitFinnhubClient() *finnhub.DefaultApiService {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", os.Getenv("FINNHUB_API_KEY"))
	return finnhub.NewAPIClient(cfg).DefaultApi
}

func symbolToFloat(ticker string) float64 {
	switch ticker {
	case "AAPL":
		return 1.0
	case "MSFT":
		return 2.0
	case "GOOG":
		return 3.0
	case "TSLA":
		return 4.0
	case "NVDA":
		return 5.0
	default:
		return 0.0
	}
}

func fetchAndInsertQuote(client *finnhub.DefaultApiService, ticker string) {
	quote, _, err := client.Quote(context.Background()).Symbol(ticker).Execute()
	now := time.Now().UTC()
	if err != nil || quote.O == nil || quote.C == nil {
		log.Printf("⚠️  %s failed: %v", ticker, err)
		return
	}

	vec := []float64{
		symbolToFloat(ticker),
		float64(now.Unix()) + float64(now.Nanosecond())/1e9,
		float64(*quote.O),
		float64(*quote.H),
		float64(*quote.L),
		float64(*quote.C),
	}

	req := insertReq{
		Values: vec,
		Metadata: map[string]string{
			"ticker": ticker,
			"time":   now.Format(time.RFC3339Nano),
		},
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post("http://localhost:3000/insert", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ Insert error for %s: %v", ticker, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("✅ %s inserted at %s", ticker, now.Format("15:04:05.000"))
}

// GetTradingDay returns a date (UTC) that represents the correct trading day
func GetTradingDay(now time.Time) time.Time {
	nyLoc, _ := time.LoadLocation("America/New_York")
	nowNY := now.In(nyLoc)

	// Fallback if before 9:30 AM
	marketOpen := time.Date(nowNY.Year(), nowNY.Month(), nowNY.Day(), 9, 30, 0, 0, nyLoc)
	if nowNY.Before(marketOpen) {
		nowNY = nowNY.AddDate(0, 0, -1)
	}

	// If weekend, fallback to Friday
	switch nowNY.Weekday() {
	case time.Saturday:
		nowNY = nowNY.AddDate(0, 0, -1)
	case time.Sunday:
		nowNY = nowNY.AddDate(0, 0, -2)
	}

	// Return zeroed-out time in UTC
	return time.Date(nowNY.Year(), nowNY.Month(), nowNY.Day(), 0, 0, 0, 0, time.UTC)
}

// FetchQuoteVector returns a vector [Open, High, Low, Close, TradingDayFloat]
func FetchQuoteVector(client *finnhub.DefaultApiService, ticker string) ([]float32, time.Time, error) {
	quote, _, err := client.Quote(context.Background()).Symbol(ticker).Execute()
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("API call failed: %w", err)
	}
	if quote.O == nil || quote.H == nil || quote.L == nil || quote.C == nil {
		return nil, time.Time{}, fmt.Errorf("incomplete quote data")
	}

	tradingDay := GetTradingDay(time.Now())
	daysSinceEpoch := float32(tradingDay.Unix()) / 86400.0

	vector := []float32{
		*quote.O,
		*quote.H,
		*quote.L,
		*quote.C,
		daysSinceEpoch,
	}

	return vector, tradingDay, nil
}

// func main() {
// 	apiKey := os.Getenv("FINNHUB_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("❌ FINNHUB_API_KEY not set.")
// 	}

// 	client := InitFinnhubClient()
// 	tickers := []string{"AAPL", "MSFT", "GOOG", "TSLA"} // Example: 4 tickers

// 	delayBetweenBatches := time.Duration(60/len(tickers)) * time.Second
// 	delayBetweenRequests := 15 * time.Millisecond

// 	for {
// 		start := time.Now()
// 		for _, ticker := range tickers {
// 			fetchAndInsertQuote(client, ticker)
// 			time.Sleep(delayBetweenRequests)
// 		}

// 		// wait remaining time until the next batch
// 		elapsed := time.Since(start)
// 		if elapsed < delayBetweenBatches {
// 			time.Sleep(delayBetweenBatches - elapsed)
// 		}
// 	}
// }