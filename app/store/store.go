package store

// see https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type UrlMapper interface {
	Write(string, string) error
	Read(string) (string, error)
}

// dynamo db implement
type DynamoDbUrlMapper struct {
	TableName string
}

// re-use aws session
var globalSess *session.Session

func genSess() (*session.Session, error) {
	if globalSess == nil {
		globalSess, err := session.NewSession(&aws.Config{
			Region: aws.String("ap-northeast-1"),
			// TODO: dev時docker network name, test時localhost, prod時?と
			// 環境変数でコントロールできる必要がある
			// 現状だと単体テストが落ちる(devを優先した)
			Endpoint: aws.String("http://dynamodb-local:8000"),
		})
		return globalSess, err
	}
	return globalSess, nil
}

func (dynamo DynamoDbUrlMapper) Write(origin string, shorten string) error {
	sess, err := genSess()
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)
	putParams := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"shorten": {
				S: aws.String(shorten),
			},
			"origin": {
				S: aws.String(origin),
			},
		},
		TableName: aws.String(dynamo.TableName),
	}

	_, err = svc.PutItem(putParams)
	if err != nil {
		return err
	}

	return nil
}

func (dynamo DynamoDbUrlMapper) Read(shorten string) (string, error) {
	sess, err := genSess()
	if err != nil {
		return "", err
	}

	svc := dynamodb.New(sess)
	getParams := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"shorten": {
				S: aws.String(shorten),
			},
		},
		TableName: aws.String(dynamo.TableName),
	}

	result, err := svc.GetItem(getParams)
	if err != nil {
		return "", err
	}

	if len(result.Item) == 0 {
		return "", nil
	}

	fetched_origin := *result.Item["origin"].S
	return fetched_origin, nil
}
