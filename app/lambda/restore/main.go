package main

import (
	"github.com/Kohei909Otsuka/simple_url_shortener/app/store"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

var dynamoStore = store.DynamoDbUrlMapper{TableName: os.Getenv("DYNAMO_TABLE")}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	origin, err := usecase.RestoreUrl(os.Getenv("BASE_URL")+req.Path, dynamoStore)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}
	if origin == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "not found",
		}, nil
	}

	headers := make(map[string]string)
	headers["Location"] = origin
	return events.APIGatewayProxyResponse{
		StatusCode: 301,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(handler)
}
