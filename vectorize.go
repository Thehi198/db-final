package main

import (
	"context"
	"fmt"
	"os"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

// Initializes the Finnhub client
func InitFinnhubClient() *finnhub.DefaultApiService {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", os.Getenv("FINNHUB_API_KEY"))
	apiClient := finnhub.NewAPIClient(cfg)
	return apiClient.DefaultApi
}

// Returns the correct trading day in UTC based on current time
func GetTradingDay(now time.Time) time.Time {
	nyLoc, _ := time.LoadLocation("America/New_York")
	nowNY := now.In(nyLoc)

	// Revert to previous day if before 9:30 AM EST
	marketOpen := time.Date(nowNY.Year(), nowNY.Month(), nowNY.Day(), 9, 30, 0, 0, nyLoc)
	if nowNY.Before(marketOpen) {
		nowNY = nowNY.AddDate(0, 0, -1)
	}

	// Revert to Friday if weekend
	switch nowNY.Weekday() {
	case time.Saturday:
		nowNY = nowNY.AddDate(0, 0, -1)
	case time.Sunday:
		nowNY = nowNY.AddDate(0, 0, -2)
	}

	// Strip time and return UTC date
	tradingDay := time.Date(nowNY.Year(), nowNY.Month(), nowNY.Day(), 0, 0, 0, 0, time.UTC)
	return tradingDay
}

// FetchQuoteVector retrieves [O, H, L, C, daysSinceEpoch] for the given ticker
func FetchQuoteVector(client *finnhub.DefaultApiService, ticker string) ([]float32, time.Time, error) {
	quote, _, err := client.Quote(context.Background()).Symbol(ticker).Execute()
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("API call failed: %w", err)
	}

	if quote.O == nil || quote.H == nil || quote.L == nil || quote.C == nil {
		return nil, time.Time{}, fmt.Errorf("missing quote data in response")
	}

	tradingDay := GetTradingDay(time.Now())
	unixSeconds := float64(tradingDay.Unix())
	daysSinceEpoch := float32(unixSeconds / 86400.0)

	vector := []float32{
		*quote.O,
		*quote.H,
		*quote.L,
		*quote.C,
		daysSinceEpoch,
	}

	return vector, tradingDay, nil
}
