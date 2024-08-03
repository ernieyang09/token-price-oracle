package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
	Host	 string
	User	 string
	Port	 int
	Password string
	Database string
}

func (m *Mysql) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        m.User, m.Password, m.Host, m.Port, m.Database)
}

func InitDB(config *Mysql) (*gorm.DB, error) {
	dsn := config.DSN()
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}