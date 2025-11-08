package coin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func RegisterRoutes(r *mux.Router) {
	// GET /api/go/coin?symbol=BTC
	r.HandleFunc("", getCoin()).Methods("GET")
}

// get Coin by Symbol from CoinMarketApi
func getCoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1) Read Symbol from Query String
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			http.Error(w, "Missing Query Param: symbol", http.StatusBadRequest)
			return
		}

		// 2) read API key, Base Url, and Path from env
		apiKey := os.Getenv("X-CMC_PRO_API_KEY")
		baseUrl := os.Getenv("CMC_BASE_URL")
		path := os.Getenv("CMC_LATEST_QUOTES_PATH")
		if apiKey == "" {
			http.Error(w, "Server Not Configured. Missing X-CMC_PRO_API_KEY", http.StatusBadRequest)
			return
		} else if baseUrl == "" {
			http.Error(w, "Server Not Configured. Missing CMC_BASE_URL", http.StatusBadRequest)
			return
		} else if path == "" {
			http.Error(w, "Server Not Configured. Missing CMC_LATEST_QUOTES_PATH", http.StatusBadRequest)
			return
		}

		// 3) Build request string
		url := fmt.Sprintf("%s%s?symbol=%s", baseUrl, path, symbol)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 4) Headers
		req.Header.Set("X-CMC_PRO_API_KEY", apiKey)
		req.Header.Set("Accept", "application/json")

		// 5) Set Timeout
		client := &http.Client{Timeout: 5 * time.Second}

		// 6) Execute Request and Close Input Stream to API
		response, err := client.Do(req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		defer response.Body.Close()

		var resultJson map[string]any
		err = json.NewDecoder(response.Body).Decode(&resultJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resultJson)
	}
}
