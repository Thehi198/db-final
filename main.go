func main() {
	apiKey := os.Getenv("POLYGON_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set your POLYGON_API_KEY environment variable.")
	}

	ticker := "AAPL"

	entry, dateStr, err := fetchLatestAvailableDay(ticker, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	dailyVec := buildDailyVector(*entry, ticker)
	fmt.Println("Embedded vector for latest available date", dateStr, ":", dailyVec)
}