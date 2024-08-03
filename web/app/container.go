package app

import (
	"fmt"
	"log"
	"oracle-go/pkg/db"

	"github.com/go-ini/ini"
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"gorm.io/gorm"
)



var container = dig.New()


func Provide(constructor interface{}, opts ...dig.ProvideOption) error {
	return container.Provide(constructor, opts...)
}

func Invoke(function interface{}, opts ...dig.InvokeOption) error {
	return container.Invoke(function, opts...)
}

func InitializeContainer() (*dig.Container) {
	
	err := Provide(func() (*ini.File, error) {
		return ini.Load("conf.ini")
	})

	if err != nil {
		log.Fatalf("Failed to provide configuration: %v", err)
	}


	err = Provide(func(cfgs *ini.File) (*gorm.DB, error) {
		mysqlSection := cfgs.Section("mysql")
		return db.InitDB(&db.Mysql{
			Host: mysqlSection.Key("host").String(),
			User: mysqlSection.Key("user").String(),
			Port: mysqlSection.Key("port").MustInt(),
			Password: mysqlSection.Key("password").String(),
			Database: mysqlSection.Key("database").String(),
		})
	})

	if err != nil {
		log.Fatalf("Failed to provide DB: %v", err)
	}

	err = Provide(func(cfgs *ini.File) (*redis.Client) {
		redisSection := cfgs.Section("redis")
		return redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", redisSection.Key("host").String(), redisSection.Key("port").String()),
			Password: redisSection.Key("password").String(),
			DB: redisSection.Key("database").MustInt(),
		})
	})

	if err != nil {
		log.Fatalf("Failed to provide Redis: %v", err)
	}

   

	return container

}