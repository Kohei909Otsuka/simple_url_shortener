.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	rm -rf ./hello-world/hello-world
	rm -rf ./app/lambda/shorten/shorten

build:
	GOOS=linux GOARCH=amd64 go build -o hello-world/hello-world ./hello-world
	GOOS=linux GOARCH=amd64 go build -o app/lambda/shorten/shorten ./app/lambda/shorten
