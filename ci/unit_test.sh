#!/bin/sh

set -eu

echo "unit test started"

# curl used to install dep, git used by dep internally
apk add curl git

# move souce to $GOPATH in go lang docker container
mkdir -p /go/src/github.com/Kohei909Otsuka/simple_url_shortener
mv gopath/src/github.com/Kohei909Otsuka/simple_url_shortener/* \
   /go/src/github.com/Kohei909Otsuka/simple_url_shortener/

# cd to repository root in GOPATH
cd /go/src/github.com/Kohei909Otsuka/simple_url_shortener

# install dep
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# install dependecy
dep ensure

# run unit test
# https://github.com/golang/go/issues/26988#issue-350494515
# store unit test fails cuz of docker is not installed in the docker container runned by ci
CGO_ENABLED=0 go test -v ./app/usecase/ ./app/entity/

echo "unit test done"
