package main

import (
	"net/http"
	"oracle-go/web/app"
	"oracle-go/web/oracle/controller"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"gorm.io/gorm"
)



func main() {
	container := app.InitializeContainer()
	
	router := gin.Default()
	v1APIGroup := router.Group("/api/v1")
	
	controller.OracleRegisterController(container, v1APIGroup)
	
	container.Invoke(func(db *gorm.DB, cfgs *ini.File) {
		router.GET("/healthy", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		port := cfgs.Section("web").Key("port").String()
		if port == "" {
			port = "8080" // Default port
		}

		router.Run(":" + port)
	})
		
	

}