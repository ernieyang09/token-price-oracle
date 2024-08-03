package strategy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oracle-go/pkg/logger"
	"time"

	"github.com/go-ini/ini"
	"github.com/robfig/cron/v3"
)

type CoinGecko struct {
	Cfg *ini.Section
	ApiKey string
}

type CoinGeckoResponse map[string]struct {
	USD float64 `json:"usd"`
}


func (cg *CoinGecko) Init(cron *cron.Cron) {
	cg.ApiKey = cg.Cfg.Key("apikey").String()
	cron.AddFunc(cg.Cfg.Key("cronExpress").String(), cg.GetPrice)
}


func (cg *CoinGecko) GetPrice() {
	logger.Info("Get price from CoinGecko")

	coinId := "ethereum" // Replace with your actual coin ID
	url := fmt.Sprintf("https://pro-api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coinId)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accepts", "application/json")
	req.Header.Add("x-cg-demo-api-key", cg.ApiKey)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API request failed with status: %s\n", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var response CoinGeckoResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	
	logger.Info(fmt.Sprintf("CoinGecko|eth|%.9f", response[coinId].USD))

}