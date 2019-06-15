#!/bin/sh

set -eu

echo "integration test started"

apk add gcc make libc-dev g++

cd sus/integration_test

bundle install

BASE_URL=$BASE_URL bundle exec rspec shorten_url_spec.rb

echo "integration test done"
