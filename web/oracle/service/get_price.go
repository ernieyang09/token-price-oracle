package service

import (
	"fmt"
	"oracle-go/pkg/entity"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Handlers struct {
	dig.In
	RedisClient *redis.Client
	DB          *gorm.DB
	Cfgs        *ini.File
}
const (
	TokenPrice = "price.value"
	TokenPriceId = "price.id"
)


func (h *Handlers) GetPrice(c *gin.Context) {
	tokenId := c.Param("tokenId")

	price := h.GetPriceFromRedis(c, tokenId)

	if price == nil {
		price = h.GetPriceFromDB(c, tokenId)

		if price == nil {
			c.JSON(404, gin.H{"error": "Price not found"})
			return
		}

		go func() {
			h.RedisClient.Set(c, fmt.Sprintf("%s.%s", TokenPrice, tokenId), *price, 15 * time.Second)
		}()
	}

	c.JSON(200, gin.H{
		"tokenId": tokenId,
		"price":   price,
	})	
}

func (h *Handlers) GetPriceFromRedis(c *gin.Context, tokenId string) *string {
	redisKey := fmt.Sprintf("%s.%s", TokenPrice, tokenId)
	price, err := h.RedisClient.Get(c, redisKey).Result()

	if err == redis.Nil {
		return nil
	}

	if err != nil {
		panic(err)
	}
	
	return &price
}

func (h *Handlers) GetPriceFromDB(c *gin.Context, tokenId string) *string {
	var record entity.Price
	err := h.DB.Where("token_id = ?", tokenId).
		Order("created_at desc").
		First(&record).Error

	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		panic(err)
	}
	priceStr := strconv.FormatFloat(float64(record.Price), 'f', -1, 32)
	return &priceStr
}