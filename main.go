package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var errorLog = log.New(os.Stderr, "ERROR ", log.Llongfile)

type contact struct {
	Channel string `json:"channel"`
	Address string `json:"address"`
}

// req = short for request
func reqTypeRouter(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return getChannel(req)
	case "POST":
		return addChannel(req)
	default:
		return clientsideErr(http.StatusMethodNotAllowed)
	}
}

func getChannel(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	channel := req.QueryStringParameters["channel"]
	// cont = short for contact
	cont, err := getChannelItem(channel)
	if err != nil {
		return serversideErr(err)
	}
	if cont == nil {
		return clientsideErr(http.StatusNotFound)
	}

	// jsonResp short for json response (marshalled)
	jsonResp, err := json.Marshal(cont)
	if err != nil {
		return serversideErr(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonResp),
	}, nil
}

func addChannel(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return clientsideErr(http.StatusNotAcceptable)
	}

	cont := new(contact)
	err := json.Unmarshal([]byte(req.Body), cont)
	if err != nil {
		return clientsideErr(http.StatusUnprocessableEntity)
	}

	err = addChannelItem(cont)
	if err != nil {
		return serversideErr(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/contact?channel=%s", cont.Channel)},
	}, nil
}

func clientsideErr(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func serversideErr(err error) (events.APIGatewayProxyResponse, error) {
	errorLog.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func main() {
	lambda.Start(reqTypeRouter)
}
