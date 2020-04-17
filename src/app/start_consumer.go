package app

import (
	"encoding/json"
	"fmt"
	"github.com/callingsid/shopping_bullwinkle/src/service"
	"github.com/callingsid/shopping_utils/domain"
	"github.com/callingsid/shopping_utils/logger"
	"github.com/callingsid/shopping_utils/queue"
	"os"
	"os/signal"
)

func startMessageConsumer() {
	topics := []string{ "shop2.operation.post", "shop2.operation.res"}
	consumer, errors := queue.Client.StartConsumer(topics)
	if errors != nil {
		fmt.Println("The error in the handle is :", errors)
		logger.Info("failed to get kafka client handle")
		//panic(errors)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				var request domain.Request
				err := json.Unmarshal(msg.Value, &request)
				if err != nil {
					logger.Error("consumer unmarshal err", err)
					//panic(err)
				}
				logger.Info(fmt.Sprintf("The Topic consumed  is %s", request.Topic))
				go processTopics(request)
			case consumerError := <-errors:
				msgCount++
				fmt.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				doneCh <- struct{}{}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
}

func processTopics(request domain.Request) {
	if request.Topic == "shop2.operation.post" {
		service.ProcessOpsRequest(request)
	}

	if request.Topic == "shop2.operation.res" {
		service.ProcessOpsResponse(request)
	}
}
