package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

// Structure for Polygon aggregates API response
type AggregatesResponse struct {
	Results []struct {
		T int64   `json:"t"` // Timestamp in milliseconds
		C float64 `json:"c"` // Close price
		V float64 `json:"v"` // Volume
		H float64 `json:"h"` // High price
		L float64 `json:"l"` // Low price
	} `json:"results"`
}

// Structure to hold daily entry
type DailyEntry struct {
	Timestamp time.Time
	Close     float32
	Volume    float64
	High      float32
	Low       float32
}

func fetchSingleDayAggregate(ticker, apiKey, date string) (*DailyEntry, error) {
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/1/day/%s/%s?adjusted=true&sort=asc&limit=1&apiKey=%s",
		ticker,
		date,
		date,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch aggregates: status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var agg AggregatesResponse
	err = json.Unmarshal(body, &agg)
	if err != nil {
		return nil, err
	}

	if len(agg.Results) == 0 {
		return nil, fmt.Errorf("no data available for date: %s", date)
	}

	res := agg.Results[0]
	entry := &DailyEntry{
		Timestamp: time.Unix(0, res.T*int64(time.Millisecond)),
		Close:     float32(res.C),
		Volume:    res.V,
		High:      float32(res.H),
		Low:       float32(res.L),
	}

	return entry, nil
}

func fetchLatestAvailableDay(ticker, apiKey string) (*DailyEntry, string, error) {
	date := time.Now().AddDate(0, 0, -1) // Start from yesterday

	for {
		dateStr := date.Format("2006-01-02")
		fmt.Println("Trying date:", dateStr)

		entry, err := fetchSingleDayAggregate(ticker, apiKey, dateStr)
		if err == nil {
			return entry, dateStr, nil
		}

		fmt.Println("No data for", dateStr, "going back one day...")
		date = date.AddDate(0, 0, -1)
	}
}

func tickerEmbedding(ticker string) []float32 {
	embed := make([]float32, len(ticker))
	for i, c := range ticker {
		ascii := int(c)
		embed[i] = float32(ascii-65) / 25.0 // normalize to [0,1]
	}
	return embed
}

func time2vec(t time.Time) []float32 {
	ts := float64(t.Unix())
	scale := 100000.0
	return []float32{
		float32(math.Sin(ts / scale)),
		float32(math.Cos(ts / scale)),
	}
}

func buildDailyVector(entry DailyEntry, ticker string) []float32 {
	tickerVec := tickerEmbedding(ticker)
	timeVec := time2vec(entry.Timestamp)
	dataVec := []float32{entry.Close, float32(entry.Volume), entry.High, entry.Low}

	vector := make([]float32, 0, len(tickerVec)+len(timeVec)+len(dataVec))
	vector = append(vector, tickerVec...)
	vector = append(vector, timeVec...)
	vector = append(vector, dataVec...)

	return vector
}
