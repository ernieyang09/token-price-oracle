package strategy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oracle-go/pkg/logger"
	"strconv"
	"time"

	"github.com/go-ini/ini"
	"github.com/robfig/cron/v3"
)

type Bitget struct {
	Cfg *ini.Section
}

type Data struct {
	AskPr float32 `json:"askPr"`
	BidPr float32 `json:"bidPr"`
}

func (d *Data) UnmarshalJSON(data []byte) error {
    var temp struct {
        BidPr string `json:"bidPr"`
        AskPr string `json:"askPr"`
    }

    if err := json.Unmarshal(data, &temp); err != nil {
        return err
    }

    bidPr, err := strconv.ParseFloat(temp.BidPr, 32)
    if err != nil {
        return err
    }

    askPr, err := strconv.ParseFloat(temp.AskPr, 32)
    if err != nil {
        return err
    }

    d.BidPr = float32(bidPr)
    d.AskPr = float32(askPr)

    return nil
}

type BitgetResponse struct {
	Data        []Data `json:"data"`
}

func (bg *Bitget) Init(cron *cron.Cron) {
	cron.AddFunc(bg.Cfg.Key("cronExpress").String(), bg.GetPrice)
}


func (bg *Bitget) GetPrice() {
	logger.Info("Get price from Bitget")

	symbol := "ETHUSDT"
	url := fmt.Sprintf("https://api.bitget.com/api/v2/spot/market/tickers?symbol=%s", symbol)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

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

	var response BitgetResponse

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	price := (response.Data[0].AskPr + response.Data[0].BidPr) / 2

	
	logger.Info(fmt.Sprintf("Biget|eth|%.9f", price))

}