package store_test

import (
	"fmt"
	"github.com/Kohei909Otsuka/simple_url_shortener/app/store"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"os/exec"
	"testing"
)

// re-use aws session
var globalSess *session.Session

func genSess() (*session.Session, error) {
	if globalSess == nil {
		globalSess, err := session.NewSession(&aws.Config{
			Region:   aws.String(os.Getenv("AWS_REGION")),
			Endpoint: aws.String(os.Getenv("AWS_ENDPOINT")),
		})
		return globalSess, err
	}
	return globalSess, nil
}

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
		TableName: aws.String(os.Getenv("DYNAMO_TABLE")),
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
		TableName: aws.String(os.Getenv("DYNAMO_TABLE")),
	}

	_, err := svc.DeleteTable(input)
	if err != nil {
		fmt.Printf("could not delete table %s", err)
	}
}

func TestMain(m *testing.M) {
	// before
	runCmd := "docker run -d -p 8000:8000 amazon/dynamodb-local"
	runOut, runErr := exec.Command("/bin/sh", "-c", runCmd).Output()
	if runErr != nil {
		fmt.Printf("failed to start docker container. err: %s", runErr)
		os.Exit(1)
	}
	containerId := string(runOut)

	os.Setenv("BASE_URL", "https://shortener.com")
	os.Setenv("AWS_REGION", "ap-northeast-1")
	os.Setenv("AWS_ENDPOINT", "http://localhost:8000")
	os.Setenv("DYNAMO_TABLE", "test_urls")

	code := m.Run()

	// after
	rmCmd := fmt.Sprintf("docker container rm -f %s", containerId)
	_, rmErr := exec.Command("/bin/sh", "-c", rmCmd).Output()
	if rmErr != nil {
		fmt.Printf("failed to rm docker container. err: %s", rmErr)
		os.Exit(1)
	}
	os.Setenv("BASE_URL", "")
	os.Setenv("AWS_REGION", "")
	os.Setenv("AWS_ENDPOINT", "")
	os.Setenv("DYNAMO_TABLE", "")
	os.Exit(code)
}

func TestDynamoDbWrite(t *testing.T) {
	createTable()

	origin := "http://original.com"
	shorten := "http://shorten.com/abcdefg"

	dynamoStore := store.DynamoDbUrlMapper{TableName: os.Getenv("DYNAMO_TABLE")}
	err := dynamoStore.Write(origin, shorten)
	if err != nil {
		t.Errorf("Dynamo Write failed, err is %s", err)
	}

	sess, _ := genSess()

	svc := dynamodb.New(sess)
	getParams := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"shorten": {
				S: aws.String(shorten),
			},
		},
		TableName: aws.String(dynamoStore.TableName),
	}

	result, err := svc.GetItem(getParams)
	if err != nil {
		t.Errorf("Dynamo Read failed, err is %s", err)
	}

	fetched_origin := *result.Item["origin"].S
	if fetched_origin != origin {
		t.Errorf("Dynamo Write failed, should write %s but write %s", origin, fetched_origin)
	}

	deleteTable()
}

func TestDynamoDbRead(t *testing.T) {
	createTable()

	origin := "http://original.com"
	shorten := "http://shorten.com/abcdefg"

	dynamoStore := store.DynamoDbUrlMapper{TableName: os.Getenv("DYNAMO_TABLE")}

	resultEmpty, err := dynamoStore.Read(shorten)
	if err != nil {
		t.Errorf("Dynamo Read faild %s", err)
	}

	if resultEmpty != "" {
		t.Errorf("Dynamo Read faild, should fetch empty but got %s", resultEmpty)
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
