package vectordb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

// Structure for Finnhub API candle response
type CandleResponse struct {
	C []float64 `json:"c"`
	V []float64 `json:"v"`
	T []int64   `json:"t"`
	S string    `json:"s"`
}

// Structure to hold a daily entry
type DailyEntry struct {
	Timestamp time.Time
	Close     float32
	Volume    float64
}

// Fetches a single day's data from Finnhub
func fetchSingleDayAggregate(ticker, apiKey, date string) (*DailyEntry, error) {
	from, _ := time.Parse("2006-01-02", date)
	to := from.AddDate(0, 0, 1)

	fromTs := from.Unix()
	toTs := to.Unix()

	url := fmt.Sprintf("https://finnhub.io/api/v1/stock/candle?symbol=%s&resolution=D&from=%d&to=%d&token=%s",
		ticker, fromTs, toTs, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch candles: status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var candle CandleResponse
	err = json.Unmarshal(body, &candle)
	if err != nil {
		return nil, err
	}

	if candle.S != "ok" || len(candle.C) == 0 {
		return nil, fmt.Errorf("no data available for date: %s", date)
	}

	entry := &DailyEntry{
		Timestamp: time.Unix(candle.T[0], 0),
		Close:     float32(candle.C[0]),
		Volume:    candle.V[0],
	}
	return entry, nil
}

// Finds the latest available day with data
func fetchLatestAvailableDay(ticker, apiKey string) (*DailyEntry, string, error) {
	date := time.Now().AddDate(0, 0, -1) // start from yesterday
	for {
		dateStr := date.Format("2006-01-02")
		entry, err := fetchSingleDayAggregate(ticker, apiKey, dateStr)
		if err == nil {
			return entry, dateStr, nil
		}
		date = date.AddDate(0, 0, -1)
	}
}

// Embeds ticker string into a vector
func tickerEmbedding(ticker string) []float32 {
	embed := make([]float32, len(ticker))
	for i, c := range ticker {
		ascii := int(c)
		embed[i] = float32(ascii-65) / 25.0
	}
	return embed
}

// Embeds timestamp into a time2vec vector
func time2vec(t time.Time) []float32 {
	ts := float64(t.Unix())
	scale := 100000.0
	return []float32{
		float32(math.Sin(ts / scale)),
		float32(math.Cos(ts / scale)),
	}
}

// Builds the final daily vector combining ticker, time, close, and volume
func buildDailyVector(entry DailyEntry, ticker string) []float32 {
	tickerVec := tickerEmbedding(ticker)
	timeVec := time2vec(entry.Timestamp)
	dataVec := []float32{entry.Close, float32(entry.Volume)}

	vector := make([]float32, 0, len(tickerVec)+len(timeVec)+len(dataVec))
	vector = append(vector, tickerVec...)
	vector = append(vector, timeVec...)
	vector = append(vector, dataVec...)

	return vector
}

// FINAL exported function: fetches and builds latest available vector
func FetchLatestVector(ticker, apiKey string) ([]float32, string, error) {
	entry, dateStr, err := fetchLatestAvailableDay(ticker, apiKey)
	if err != nil {
		return nil, "", err
	}
	vec := buildDailyVector(*entry, ticker)
	return vec, dateStr, nil
}
