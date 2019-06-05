# 開発

dynamo db localの起動をネットワーク指定して、
lambdaが動いているcontainerから見えるようにする

```
docker network create aws-local
docker-compose up
sam local start-api --docker-network simple_url_shortener_aws-local
```

# テスト

単体テストすべて実行
```
go test ./...
```

# スクリプト

shoten api

```
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"origin":"https://origin.com"}' \
  localhost:3000
```

dynamo

```
aws dynamodb list-tables --endpoint-url http://localhost:8000

aws dynamodb create-table --table-name 'urls' \
  --attribute-definitions '[{"AttributeName":"shorten","AttributeType": "S"}]' \
  --key-schema '[{"AttributeName":"shorten","KeyType": "HASH"}]' \
  --provisioned-throughput '{"ReadCapacityUnits": 5,"WriteCapacityUnits": 5}' \
  --endpoint-url http://localhost:8000
  ```
