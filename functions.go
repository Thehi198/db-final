package main

import (
	"context"
	"fmt"
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
	apiClient := finnhub.NewAPIClient(cfg)
	return apiClient.DefaultApi
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
	default:
		return 0.0
	}
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
