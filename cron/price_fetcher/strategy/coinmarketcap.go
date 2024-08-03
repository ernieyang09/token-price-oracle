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

type CoinMartketCap struct {
	Cfg *ini.Section
	ApiKey string
}


type CoinMarketCapResponse struct {
    Data map[string]struct {
        Quote struct {
            USD struct {
                Price float64 `json:"price"`
            } `json:"USD"`
        } `json:"quote"`
    } `json:"data"`
}

func (cmc *CoinMartketCap) Init(cron *cron.Cron) {
	cmc.ApiKey = cmc.Cfg.Key("apikey").String()
	cron.AddFunc(cmc.Cfg.Key("cronExpress").String(), cmc.GetPrice)
}

func (cmc *CoinMartketCap) GetPrice() {
	logger.Info("Get price from CoinMarketCap")

	coinId := "1027"
	url := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?id=%s", coinId)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", cmc.ApiKey)

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
	}

	var response CoinMarketCapResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	logger.Info(fmt.Sprintf("CoinMarketCap|eth|%.9f", response.Data[coinId].Quote.USD.Price))
}