#!/bin/sh

# NOTE: expecting compile is already done on local

set -eu

echo "deploy started"

apk add gcc musl-dev jq

# export PATH for installed python scritp
export PATH=$PATH:/root/.local/bin

# install aws cli
pip3 install awscli --user

# install sam cli
pip3 install aws-sam-cli --user

# package
sam package \
  --template-file ./sus/template.yaml \
  --s3-bucket simple-url-shortener \
  --output-template-file ./sus/packaged.yaml

# deploy
sam deploy \
  --template-file ./sus/packaged.yaml \
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
while read api_id
do
  aws apigateway update-stage \
    --rest-api-id <api-id> \
    --stage-name Prod \
    --patch-operations \
      op=replace,path=/*/*/logging/dataTrace,value=true \
      op=replace,path=/*/*/logging/loglevel,value=Info \
      op=replace,path=/*/*/metrics/
done < rest-api_id.txt

echo "deploy done"
