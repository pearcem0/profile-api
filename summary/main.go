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

var stage = getEnv("STAGE", "dev")
var region = getEnv("REGION", "eu-west-2")
var errorLog = log.New(os.Stderr, "ERROR ", log.Llongfile)

type summary struct {
	Section string `json:"section"`
	Content string `json:"content"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// req = short for request
func reqTypeRouter(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return getSection(req)
	case "POST":
		return addSection(req)
	default:
		return clientsideErr(http.StatusMethodNotAllowed)
	}
}

func getSection(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	section := req.QueryStringParameters["section"]
	// summ = short for summary
	summ, err := getSectionItem(section)
	if err != nil {
		return serversideErr(err)
	}
	if summ == nil {
		return clientsideErr(http.StatusNotFound)
	}

	// jsonResp short for json response (marshalled)
	jsonResp, err := json.Marshal(summ)
	if err != nil {
		return serversideErr(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonResp),
	}, nil
}

func addSection(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return clientsideErr(http.StatusNotAcceptable)
	}

	summ := new(summary)
	err := json.Unmarshal([]byte(req.Body), summ)
	if err != nil {
		return clientsideErr(http.StatusUnprocessableEntity)
	}

	err = addSectionItem(summ)
	if err != nil {
		return serversideErr(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/summary?section=%s", summ.Section)},
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
