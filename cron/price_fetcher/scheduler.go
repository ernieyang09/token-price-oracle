package main

import (
	"oracle-go/cron/price_fetcher/strategy"
	"oracle-go/pkg/env"
	"sync"

	"github.com/go-ini/ini"
	cron "github.com/robfig/cron/v3"
)

func SetupJobs(c *cron.Cron) {
	wg := sync.WaitGroup{}

	configName := env.GetConfigValue("CONFIG_NAME", "conf.ini")

	cfgs, err := ini.Load(configName)

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	// for each strategy in folder, create a job
	jobs := []strategy.GetPriceStrategy{
		&strategy.CoinGecko{Cfg: cfgs.Section("coingecko")},
		&strategy.CoinMartketCap{Cfg: cfgs.Section("coinmarketcap")},
		&strategy.Bitget{Cfg: cfgs.Section("bitget")},
	}

	for _, job := range jobs {
		wg.Add(1)

		job.Init(c)
		go func(j strategy.GetPriceStrategy) {
			defer wg.Done()
			j.GetPrice()
		}(job)

		wg.Wait()
	}
}