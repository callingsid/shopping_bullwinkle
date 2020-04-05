package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/callingsid/shopping_bullwinkle/src/service"
	"github.com/callingsid/shopping_bullwinkle/src/utils"
	"github.com/callingsid/shopping_utils/logger"
	"github.com/callingsid/shopping_utils/queue"
	"github.com/callingsid/shopping_utils/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	DEADLINE               time.Duration = 10 * time.Second
	topic_ops_req			= "shop2.operation.post"
)

var httpDurationHistogram *prometheus.HistogramVec

func init() {
	logger.Info("Initializing histogram and summary metrics ...")
	httpDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time taken to process http request in seconds",
		Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5},
	}, []string{"function"})
	prometheus.MustRegister(httpDurationHistogram)
}
func getItemId(itemIdParam string) (int64, rest_errors.RestErr) {
	itemId, userErr := strconv.ParseInt(itemIdParam, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError("user id should be a number")
	}
	return itemId, nil
}

func GetRequest(c *gin.Context) {
	//lets fetch metrics
	start := time.Now()
	defer func() {
		httpDuration := time.Since(start)
		httpDurationHistogram.WithLabelValues("Shopping-Get-Request").Observe(httpDuration.Seconds())
	}()

	microservice := utils.GetFinalPathString(c.FullPath(), 2)
	logger.Info(fmt.Sprintf("The last path is %s ", microservice))

	itemId := c.Param("item_id")

	//generate uid using random string
	rand.Seed(time.Now().UnixNano())
	id := utils.RandomString(10)
	logger.Info(fmt.Sprintf("The value of uid generated in controller is : %s", id))

	//get CRUD HTTP method
	method := c.Request.Method

	kdata := make(map[string]interface{})
	kdata["data"] = itemId
	kdata["uid"] = id
	kdata["topic"] = topic_ops_req
	kdata["fwd_topic"] = microservice
	kdata["method"] = method

	if err :=  queue.PClient.Publish(topic_ops_req, kdata); err != nil {
		logger.Error("Kafka error: %s\n", err)
		c.JSON(http.StatusInternalServerError, err)
	}
	time.Sleep( 1 * time.Second)
	result, restErr := service.GetResponse(id); if restErr != nil {
		logger.Error("Http error: %s\n", restErr)
		c.JSON(restErr.Status(), restErr)
		c.Abort()
	} else {
		logger.Info(fmt.Sprintf("*********The final result returned is %s***********", result))
		c.JSON(http.StatusOK, result)
	}
}

func PostRequest(c *gin.Context) {

	//lets fetch metrics
	start := time.Now()
	defer func() {
		httpDuration := time.Since(start)
		httpDurationHistogram.WithLabelValues("Shopping-Post-Request").Observe(httpDuration.Seconds())
	}()



	//fetch the microservice name and start the microservice
	microservice := utils.GetFinalPathString(c.FullPath(), 1)
	logger.Info(fmt.Sprintf("The last path is %s ", microservice))

	//generate uid using random string
	rand.Seed(time.Now().UnixNano())
	id := utils.RandomString(10)
	logger.Info(fmt.Sprintf("The value of uid generated in controller is : %s", id))

	//get CRUD HTTP method
	method := c.Request.Method

	//setting contexts with uid and http method
	deadline := time.Now().Add(DEADLINE)
	ctx, cancel:= context.WithDeadline(context.Background(), deadline)
	defer cancel()
	ctx = utils.SetContext(ctx, id, method)

	var _req map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&_req)
	if err != nil {
		logger.Error("failed to read request body", err)
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
	}
	defer c.Request.Body.Close()

	// add data to send to publish to topic
	kdata := make(map[string]interface{})
	kdata["data"] = _req
	kdata["uid"] = id
	kdata["topic"] = topic_ops_req
	kdata["fwd_topic"] = microservice
	kdata["method"] = method

	if err :=  queue.PClient.Publish(topic_ops_req, kdata); err != nil {
		logger.Error("Kafka error: %s\n", err)
		c.JSON(http.StatusInternalServerError, err)
		c.Abort()
	}
	time.Sleep( 1 * time.Second)
	result, restErr := service.GetResponse(id); if restErr != nil {
		logger.Error("Http error: %s\n", restErr)
		c.JSON(restErr.Status(), restErr)
		c.Abort()
	}
	logger.Info(fmt.Sprintf("*********The final result returned is %s***********", result))

	c.JSON(http.StatusCreated, result)
}



