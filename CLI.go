// package main

// import (
// 	"log"
// 	"os"
// 	"time"
// )

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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"vectordb/vectordb" // assuming your module name is `db-final`

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

var db *vectordb.Database

func symbolToFloat(ticker string) float64 {
	switch strings.ToUpper(ticker) {
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

// API data
func InitFinnhubClient() *finnhub.DefaultApiService {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", os.Getenv("d09ffu9r01qnv9cikdjgd09ffu9r01qnv9cikdk0"))
	return finnhub.NewAPIClient(cfg).DefaultApi
}

func fetchAndInsertQuote(client *finnhub.DefaultApiService, ticker string, api string) {
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

	req := map[string]interface{}{
		"values": vec,
		"metadata": map[string]string{
			"ticker": ticker,
			"time":   now.Format(time.RFC3339Nano),
		},
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post(api+"/insert", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ Insert error for %s: %v", ticker, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("✅ %s inserted at %s", ticker, now.Format("15:04:05.000"))
}

// CLI commands
func runCreate(args []string, api string) {
	if len(args) < 1 {
		fmt.Println("Usage: create <dimension>")
		return
	}
	dim, _ := strconv.Atoi(args[0])
	if api == "" {
		db = vectordb.NewDatabase(dim)
		fmt.Println("✅ Local DB created with dimension", dim)
	} else {
		body := map[string]interface{}{"dimension": dim}
		sendPost(api+"/create", body)
	}
}

func runInsert(args []string, api string) {
	vec, meta := parseVectorAndFlags(args)
	if api == "" {
		if db == nil {
			fmt.Println("❌ Local DB not initialized.")
			return
		}
		id, err := db.Insert(vectordb.Vector{Values: vec, Metadata: meta})
		if err != nil {
			fmt.Println("❌ Insert error:", err)
			return
		}
		fmt.Println("✅ Inserted locally with UUID:", id)
	} else {
		body := map[string]interface{}{"values": vec, "metadata": meta}
		sendPost(api+"/insert", body)
	}
}

func runQuery(args []string, api string) {
	vec, meta, k := parseVectorQueryFlags(args)
	if api == "" {
		if db == nil {
			fmt.Println("❌ Local DB not initialized.")
			return
		}
		res, err := db.Query(vec, k, meta)
		if err != nil {
			fmt.Println("❌ Query error:", err)
			return
		}
		for _, r := range res {
			fmt.Printf("UUID: %s, Dist: %.4f, Meta: %v\n", r.Vector.UUID, r.Distance, r.Vector.Metadata)
		}
	} else {
		body := map[string]interface{}{"values": vec, "k": k, "metadata_filter": meta}
		sendPost(api+"/query", body)
	}
}

func runList(api string) {
	if api == "" {
		if db == nil {
			fmt.Println("❌ Local DB not initialized.")
			return
		}
		for _, v := range db.Vectors {
			fmt.Printf("UUID: %s, Values: %v, Meta: %v\n", v.UUID, v.Values, v.Metadata)
		}
	} else {
		resp, err := http.Get(api + "/list")
		if err != nil {
			fmt.Println("❌ List error:", err)
			return
		}
		defer resp.Body.Close()
		io.Copy(os.Stdout, resp.Body)
	}
}

func runFetchInsert(args []string, api string) {
	if len(args) < 2 || args[0] != "--ticker" || api == "" {
		fmt.Println("Usage: fetch-insert --ticker SYMBOL --api URL")
		return
	}
	ticker := args[1]
	client := InitFinnhubClient()
	fetchAndInsertQuote(client, ticker, api)
}

func runFetchLoop(api string) {
	if api == "" {
		fmt.Println("❌ Must specify --api to use fetch-loop")
		return
	}
	client := InitFinnhubClient()
	tickers := []string{"AAPL", "MSFT", "GOOG", "TSLA"}
	delayBetweenBatches := time.Duration(60/len(tickers)) * time.Second
	delayBetweenRequests := 15 * time.Millisecond

	for {
		start := time.Now()
		for _, ticker := range tickers {
			fetchAndInsertQuote(client, ticker, api)
			time.Sleep(delayBetweenRequests)
		}
		elapsed := time.Since(start)
		if elapsed < delayBetweenBatches {
			time.Sleep(delayBetweenBatches - elapsed)
		}
	}
}

func parseVectorAndFlags(args []string) ([]float64, map[string]string) {
	vec := []float64{}
	meta := map[string]string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--ticker":
			meta["ticker"] = args[i+1]
			i++
		default:
			v, _ := strconv.ParseFloat(args[i], 64)
			vec = append(vec, v)
		}
	}
	return vec, meta
}

func parseVectorQueryFlags(args []string) ([]float64, map[string]string, int) {
	vec := []float64{}
	meta := map[string]string{}
	k := 1
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--ticker":
			meta["ticker"] = args[i+1]
			i++
		case "--k":
			k, _ = strconv.Atoi(args[i+1])
			i++
		default:
			v, _ := strconv.ParseFloat(args[i], 64)
			vec = append(vec, v)
		}
	}
	return vec, meta, k
}

func sendPost(url string, body map[string]interface{}) {
	data, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("❌ HTTP error:", err)
		return
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: [create|insert|query|list|fetch-insert|fetch-loop] ... [--api URL]")
		return
	}

	apiURL := ""
	args := []string{}
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--api" && i+1 < len(os.Args) {
			apiURL = os.Args[i+1]
			i++
		} else {
			args = append(args, os.Args[i])
		}
	}

	if len(args) == 0 {
		fmt.Println("Missing command.")
		return
	}

	switch args[0] {
	case "create":
		runCreate(args[1:], apiURL)
	case "insert":
		runInsert(args[1:], apiURL)
	case "query":
		runQuery(args[1:], apiURL)
	case "list":
		runList(apiURL)
	case "fetch-insert":
		runFetchInsert(args[1:], apiURL)
	case "fetch-loop":
		runFetchLoop(apiURL)
	default:
		fmt.Println("Unknown command:", args[0])
	}
}
