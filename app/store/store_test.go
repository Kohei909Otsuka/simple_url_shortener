package store_test

import (
	"fmt"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/entity"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/store"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"testing"
)

// re-use aws session
var globalSess *session.Session

func genSess() (*session.Session, error) {
	if globalSess == nil {
		globalSess, err := session.NewSession(&aws.Config{
			Region:   aws.String("ap-northeast-1"),
			Endpoint: aws.String("http://localhost:8000"),
		})
		return globalSess, err
	}
	return globalSess, nil
}

// expecting dynamo db local is running
// see https://hub.docker.com/r/amazon/dynamodb-local
func createTable() {
	sess, _ := genSess()
	svc := dynamodb.New(sess)
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("shorten"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("shorten"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String("test_urls"),
	}
	_, err := svc.CreateTable(input)
	if err != nil {
		fmt.Printf("could not create table %s", err)
	}
}

func deleteTable() {
	sess, _ := genSess()
	svc := dynamodb.New(sess)
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String("test_urls"),
	}

	_, err := svc.DeleteTable(input)
	if err != nil {
		fmt.Printf("could not delete table %s", err)
	}
}

func TestMain(m *testing.M) {
	// before
	os.Setenv("BASE_URL", "https://shortener.com")

	code := m.Run()

	// after
	os.Setenv("BASE_URL", "")
	os.Exit(code)
}

func TestDynamoDbWrite(t *testing.T) {
	createTable()

	origin := entity.OriginalUrl("http://original.com")
	shorten := entity.ShortenUrl("http://shorten.com/abcdefg")

	dynamoStore := store.DynamoDbUrlMapper{TableName: "test_urls"}
	err := dynamoStore.Write(origin, shorten)
	if err != nil {
		t.Errorf("Dynamo Write failed, err is %s", err)
	}

	sess, _ := genSess()

	svc := dynamodb.New(sess)
	getParams := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"shorten": {
				S: aws.String(string(shorten)),
			},
		},
		TableName: aws.String(dynamoStore.TableName),
	}

	result, err := svc.GetItem(getParams)
	if err != nil {
		t.Errorf("Dynamo Read failed, err is %s", err)
	}

	fetched_origin := *result.Item["origin"].S
	if fetched_origin != string(origin) {
		t.Errorf("Dynamo Write failed, should write %s but write %s", string(origin), fetched_origin)
	}

	deleteTable()
}

func TestDynamoDbRead(t *testing.T) {
	createTable()

	origin := entity.OriginalUrl("http://original.com")
	shorten := entity.ShortenUrl("http://shorten.com/abcdefg")

	dynamoStore := store.DynamoDbUrlMapper{TableName: "test_urls"}

	resultEmpty, err := dynamoStore.Read(shorten)
	if err != nil {
		t.Errorf("Dynamo Read faild %s", err)
	}

	if string(resultEmpty) != "" {
		t.Errorf("Dynamo Read faild, should fetch empty but got %s", string(resultEmpty))
	}

	dynamoStore.Write(origin, shorten)
	result, err := dynamoStore.Read(shorten)
	if err != nil {
		t.Errorf("Dynamo Read faild %s", err)
	}

	if result != origin {
		t.Errorf("Dynamo Read faild, fetch expected %s but got %s", origin, result)
	}

	deleteTable()
}
