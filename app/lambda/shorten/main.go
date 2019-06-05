package main

import (
	"encoding/json"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/store"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

var dynamoStore = store.DynamoDbUrlMapper{TableName: os.Getenv("DYNAMO_TABLE")}

type ShortenReq struct {
	Origin string `json:"origin"`
}
type ShortenRes struct {
	Shorten string `json:"shorten"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var shortenReq ShortenReq
	err := json.Unmarshal([]byte(request.Body), &shortenReq)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "could not parse json",
		}, nil
	}

	shorten, err := usecase.ShortenUrl(shortenReq.Origin, dynamoStore)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	shortenRes := ShortenRes{shorten}
	bytes, err := json.Marshal(shortenRes)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bytes),
	}, nil
}

func main() {
	lambda.Start(handler)
}
