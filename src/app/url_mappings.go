package app

import (
	"github.com/callingsid/shopping_bullwinkle/src/controller"
	"github.com/callingsid/shopping_bullwinkle/src/controller/ping"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func mapUrls() {
	router.POST("/shopping/items", controller.PostRequest)
	router.GET("/shopping/items/:item_id", controller.GetRequest)
	router.GET("/ping", ping.Ping)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
