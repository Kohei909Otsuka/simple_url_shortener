# simple_url_shortener

It shows file struct below.

```bash
.
|-- app
    |-- entity/                 <-- entity(from Clean-Architecture)
    |-- lambda/                 <-- store aws lambda functions
    |-- store/                  <-- store implementation to save urls
    |-- usecase/                <-- use case(from Clean-Architecture)
|-- Gopkg.lock                  <-- dependecy lock file by dep
|-- Gopkg.toml                  <-- dependency file by dep
|-- Makefile                    <-- Make to automate build
|-- README.md                   <-- This file
|-- dev_env.json                <-- json to store env vars for dev environment
|-- docker-compose.yml          <-- docker-compose.yml to run aws dynamodb locally
|-- template.yaml               <-- SAM template
```

## Requirements

* [AWS CLI](https://aws.amazon.com/jp/cli/)
* [SAM CLI](https://github.com/awslabs/aws-sam-cli)
* [Docker installed](https://www.docker.com/community-edition)
* [Golang version 1.x](https://golang.org)
* [Deb](https://github.com/golang/dep)

## Setup process

### Installing dependencies

we use [deb](https://github.com/golang/dep) for Go lang third party library management.
Following command will install all we need.

```shell
dep ensure
```

### Building

Build task is defined with Makefile.
Following command will build each lambda function written in Go.

```shell
make build
```

### Local development

There are one important point to consider when develop serverless application.

Basically, AWS Lambda works like glue between AWS resources.

Lambda is function in math, whose input is other AWS resource(Event), and output is other AWS resource too.

The Problem is How can we run AWS resource(S3, Api Gateway, DynamoDB, SNS, SQS..etc) on our local machine?

Fortunately some of them like S3, Api Gateway, DynamoDB have emulator, but some of them not.

This app only uses Api Gateway, Lambda, DynamoDB, so that we can develop on local machine.

However, when we wants to use other AWS resource which has no emulator, it's time to give up local development.

```bash
# create named docker network
# lambda and dynamodb are supposed on this named network later
docker network create simple_url_shortener_aws-local

# run docker container for dynamodb
docker-compose up

# create dynamodb table for dev env
# the data saved in dynamodb will lost atter stop container
# to prevent lost, u need to use docker volume
aws dynamodb create-table --table-name 'dev_urls' \
  --attribute-definitions '[{"AttributeName":"shorten","AttributeType": "S"}]' \
  --key-schema '[{"AttributeName":"shorten","KeyType": "HASH"}]' \
  --provisioned-throughput '{"ReadCapacityUnits": 5,"WriteCapacityUnits": 5}' \
  --endpoint-url http://localhost:8000

# start api gateway locally on docker network created before
sam local start-api --docker-network simple_url_shortener_aws-local --env-vars dev_env.json
```

```bash
# call shorten api
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"origin":"https://origin.com"}' \
  localhost:3000
# => api respone will be like this {"shoten": "localhost:3000/abcdef"}

# call restore api
curl localhost:3000/abcdef
# => api restore status will be 301 redirect, to https://origin.com
```

## Packaging and deployment


First and foremost, we need a `S3 bucket` where we can upload our Lambda functions packaged as ZIP before we deploy anything - If you don't have a S3 bucket to store code artifacts then this is a good time to create one:

```bash
aws s3 mb s3://BUCKET_NAME
```

Next, run the following command to package our Lambda function to S3:

```bash
# NOTE: assuming u has aws profile named serverless which has enough permission to deploy your aws resources.
sam package \
  --template-file template.yaml \
  --s3-bucket <bucket> \
  --output-template-file packaged.yaml \
  --profile serverless
```

Next, the following command will create a Cloudformation Stack and deploy your SAM resources.

```bash
# NOTE: assuming u has aws profile named serverless which has enough permission to deploy your aws resources.
sam deploy \
  --template-file ./packaged.yaml \
  --stack-name <stack> \
  --capabilities CAPABILITY_IAM \
  --profile serverless
```

After deployment is complete you can run the following command to retrieve the API Gateway Endpoint URL:

```bash
aws cloudformation describe-stacks \
    --stack-name simple_url_shortener \
    --query 'Stacks[].Outputs'
```

## Testing

### unit test

We use `testing` package that is built-in in Golang and you can simply run the following command to run our tests:

```shell
go test -v ./app/...
```

### integrate test

Iintegrate test is important for serverless app, but not yet written.

## About custom domain

We have two choice about domain

1. use custom domain
2. use aws gateway auto generate domain

### 1. use custom domain

I personally bought domain `kho21.com`.
Then config AWS api gateway and AWS route 53 to use it.
Plz [read](https://docs.aws.amazon.com/ja_jp/apigateway/latest/developerguide/how-to-custom-domains.html) this if u want use custom domain

### 2. use aws gateway auto generate domain

After deploy, AWS api gateway will generate url by stage like `https://{some}.execute-api.${region}.amazonaws.com/Prod`
U can then config lambda env var BASE_URL from AWS console.

## Ref

- [aws lambda](https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/welcome.html)
- [aws sam doc on githbub](https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md)
- [aws sam doc](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html)
- [aws sam cli doc](https://github.com/awslabs/aws-sam-cli/tree/develop/docs)
