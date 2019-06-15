#!/bin/sh

set -eu

# package
sam package \
  --template-file template.yaml \
  --s3-bucket simple-url-shortener \
  --output-template-file packaged.yaml

# deploy
sam deploy \
  --template-file packaged.yaml \
  --stack-name simple-url-shortener \
  --capabilities CAPABILITY_IAM \
  --region ap-northeast-1

aws cloudformation list-stack-resources \
  --stack-name simple-url-shortener \
  --region ap-northeast-1 \
  | jq -r '.StackResourceSummaries[] | select(.ResourceType == "AWS::ApiGateway::RestApi") | .PhysicalResourceId' > rest-api-id.txt

# enable api gateway log
# https://github.com/awslabs/serverless-application-model/blob/master/docs/faq.rst#how-to-enable-api-gateway-logs
# maybe want to use https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#awsserverlessapi later
while read api-id
do
  aws apigateway update-stage \
    --rest-api-id <api-id> \
    --stage-name Prod \
    --patch-operations \
      op=replace,path=/*/*/logging/dataTrace,value=true \
      op=replace,path=/*/*/logging/loglevel,value=Info \
      op=replace,path=/*/*/metrics/
done < rest-api-id.txt
