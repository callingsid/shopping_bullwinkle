package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/callingsid/shopping_utils/domain"
	"github.com/callingsid/shopping_utils/logger"
	"github.com/callingsid/shopping_utils/queue"
	"github.com/callingsid/shopping_utils/rest_errors"
)

var (
	rchans 	= make(map[string] chan json.RawMessage)
	response = make(map[string] json.RawMessage)
)



func ProcessOpsRequest(request domain.Request) rest_errors.RestErr{
	logger.Info(fmt.Sprintf("**************message data after unmarshal is %s and method is %s*********:", request.UID, request.Data))

	in := make(chan json.RawMessage)
	//out := make(chan json.RawMessage)
	rchans[request.UID] = in

	// Publish to topic based on <url>/value
	logger.Info(fmt.Sprintf("publishing to topic %s", request.FwdTopic))
	if err :=  queue.PClient.Publish(request.FwdTopic, request); err != nil {
		logger.Error("Kafka error: %s\n", err)
		return err
	}
	logger.Info(fmt.Sprintf("Successfully logged to topic %s", request.FwdTopic))
	//Now wait till we get response from channel from out channel
	 waitResp(request.UID, in)
	logger.Info("Ending the process() function and posting out data")

	return nil
}

func waitResp(uid string,in chan json.RawMessage)   {
	logger.Info(fmt.Sprintf("entered waitResp method for uid %s : ", uid))
	//out := make(chan json.RawMessage)
	select {
	case r := <-in:
		// append response
		logger.Info("waiting here in waitresp for the channel")
		response[uid] = r
		logger.Info(fmt.Sprintf("The response before is %s", r))
		delete(rchans, uid)
		//return responses
	}

}

func ProcessOpsResponse(request domain.Request ) {
	if c, ok := rchans[request.UID]; ok {
		logger.Info(fmt.Sprintf("The request id is %s and data is %s", request.UID, request.Data))
		c <- request.Data
	}
}

func GetResponse(uid string) (json.RawMessage, rest_errors.RestErr) {
	//var res interface{}
	res := make(map[string]interface{})
	for k,v := range response {
		if k == uid {
			 err := json.Unmarshal([]byte(v), &res)
			 if err != nil {
			 	//TODO
			 }

			 val, ok := res["status"].(float64)
			 if !ok {
			 	//TODO
			 }
			 //check if the response is having err -- status code > 200. If yes, build the error and return
			 if val > 201 {
			 	err := rest_errors.NewRestError(res["message"].(string), int(res["status"].(float64)), res["error"].(string), nil)
			 	delete(res, k)
			 	return nil, err
			 }
			 // Else there is no error, just response
			delete(res, k)
			return v, nil
		}
	}
	logger.Info("No response in getReponse method")
	return nil, rest_errors.NewInternalServerError("error when trying to save item", errors.New("http error"))
}






