.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	rm -rf ./app/lambda/shorten/shorten
	rm -rf ./app/lambda/restore/restore

# https://github.com/docker-library/golang/issues/152
# Did not work well when build on go alpine image by ci
build:
	GOOS=linux GOARCH=amd64 go build -ldflags '-d' -tags netgo -installsuffix netgo -o app/lambda/shorten/shorten ./app/lambda/shorten
	GOOS=linux GOARCH=amd64 go build -ldflags '-d' -tags netgo -installsuffix netgo -o app/lambda/restore/restore ./app/lambda/restore
