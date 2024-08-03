package entity

type Price struct {
    ID        uint64    `gorm:"primaryKey;autoIncrement;column:id"` // `bigint unsigned` maps to `uint64`
    TokenID   string    `gorm:"size:10;not null;index"`            // `varchar(10)` maps to `string`
    Price     float64   `gorm:"not null;type:decimal(20,10)"`      // `decimal(20,10)` maps to `float64`
    CreatedAt int64     `gorm:"not null;default:UNIX_TIMESTAMP()"`  // `int` maps to `int64`
}

func (Price) TableName() string {
    return "price"
}