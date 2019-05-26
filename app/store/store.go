package store

// see https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/

import (
	"github.com/Kohei909Otsuka/simple_url_shortener/app/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type UrlMapper interface {
	Write(entity.OriginalUrl, entity.ShortenUrl) error
	Read(entity.ShortenUrl) (entity.OriginalUrl, error)
}

// dynamo db implement
type DynamoDbUrlMapper struct {
	TableName string
}

func (dynamo DynamoDbUrlMapper) Write(origin entity.OriginalUrl, shorten entity.ShortenUrl) error {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)
	putParams := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"shorten": {
				S: aws.String(string(shorten)),
			},
			"origin": {
				S: aws.String(string(origin)),
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

func (dynamo *DynamoDbUrlMapper) Read(shorten entity.ShortenUrl) (entity.OriginalUrl, error) {
	zeroValue := entity.OriginalUrl("")
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		return zeroValue, err
	}

	svc := dynamodb.New(sess)
	getParams := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"shorten": {
				S: aws.String(string(shorten)),
			},
		},
		TableName: aws.String(dynamo.TableName),
	}

	result, err := svc.GetItem(getParams)
	if err != nil {
		return zeroValue, err
	}

	if len(result.Item) == 0 {
		return zeroValue, nil
	}

	fetched_origin := *result.Item["origin"].S
	return entity.OriginalUrl(fetched_origin), nil
}
