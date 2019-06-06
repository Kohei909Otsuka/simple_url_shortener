.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	rm -rf ./hello-world/hello-world
	rm -rf ./app/lambda/shorten/shorten
	rm -rf ./app/lambda/restore/restore

build:
	GOOS=linux GOARCH=amd64 go build -o hello-world/hello-world ./hello-world
	GOOS=linux GOARCH=amd64 go build -o app/lambda/shorten/shorten ./app/lambda/shorten
	GOOS=linux GOARCH=amd64 go build -o app/lambda/restore/restore ./app/lambda/restore
