package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type CoinGeckoResponse struct {
	Prices [][]float64 `json:"prices"`
}

type PriceResponse struct {
	MainCurrency string  `json:"mainCurrency"`
	VsCurrency   string  `json:"vsCurrency"`
	OpenPrice    float64 `json:"openPrice"`
	ClosePrice   float64 `json:"closePrice"`
}

func fetchBitcoinPriceHistory(mainCurrency, vsCurrency string, days int) (float64, float64, error) {
	from := time.Now().AddDate(0, 0, -days).Unix()
	to := time.Now().Unix()
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart/range?vs_currency=%s&from=%d&to=%d", mainCurrency, vsCurrency, from, to)

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	var data CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if len(data.Prices) == 0 {
		return 0, 0, fmt.Errorf("no price data available")
	}

	openPrice, closePrice := data.Prices[0][1], data.Prices[len(data.Prices)-1][1]
	return openPrice, closePrice, nil
}

func priceHandler(w http.ResponseWriter, r *http.Request) {
	mainCurrency := os.Getenv("MAIN_CURRENCY")
	if mainCurrency == "" {
		mainCurrency = "bitcoin" // Default currency
	}
	vsCurrency := os.Getenv("VS_CURRENCY")
	if vsCurrency == "" {
		vsCurrency = "cny" // Default currency
	}

	openPrice, closePrice, err := fetchBitcoinPriceHistory(mainCurrency, vsCurrency, 2)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	priceResponse := PriceResponse{
		MainCurrency: mainCurrency,
		VsCurrency:   vsCurrency,
		OpenPrice:    openPrice,
		ClosePrice:   closePrice,
	}
	// proper header
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(priceResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", priceHandler)
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
