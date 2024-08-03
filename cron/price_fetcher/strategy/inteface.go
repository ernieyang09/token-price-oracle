package strategy

import "github.com/robfig/cron/v3"


type GetPriceStrategy interface {
	Init(cron *cron.Cron)
	GetPrice()
}