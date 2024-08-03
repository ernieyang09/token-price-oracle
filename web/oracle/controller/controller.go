package controller

import (
	"fmt"
	oracleService "oracle-go/web/oracle/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func OracleRegisterController(container *dig.Container, g *gin.RouterGroup) {

	err := container.Invoke(func(h oracleService.Handlers) {
		g.GET("/price/:tokenId", h.GetPrice)
	})

	if err != nil {
		fmt.Printf("Error during controller registration: %v\n", err)
	}

}
