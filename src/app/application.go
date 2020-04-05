package app

import (
	"github.com/callingsid/shopping_utils/logger"
	"github.com/gin-gonic/gin"
)
var (
	router = gin.Default()
)


func StartApplication() {

	go startMessageConsumer()
	logger.Info("about to start the application... ")
	mapUrls()
	router.Run(":8080")
}


