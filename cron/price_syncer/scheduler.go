package main

import (
	"oracle-go/cron/price_syncer/job"
	"oracle-go/pkg/db"
	"oracle-go/pkg/env"

	"github.com/go-ini/ini"
	cron "github.com/robfig/cron/v3"
)

func SetupJobs(c *cron.Cron) {
	// Build the path to config.ini
	configName := env.GetConfigValue("CONFIG_NAME", "conf.ini")

	cfgs, err := ini.Load(configName)

	if err != nil {
		panic(err)
	}

	mysqlSection := cfgs.Section("mysql")

	mysql, err := db.InitDB(&db.Mysql{
		Host:     mysqlSection.Key("host").String(),
		User:     mysqlSection.Key("user").String(),
		Port:     mysqlSection.Key("port").MustInt(),
		Password: mysqlSection.Key("password").String(),
		Database: mysqlSection.Key("database").String(),
	})

	if err != nil {
		panic(err)
	}

	job := job.Job{Cfg: cfgs.Section("price_syncer"), DB: mysql}

	job.Init(c)
	job.SyncPrice()
}
