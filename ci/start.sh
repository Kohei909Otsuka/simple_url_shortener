#!/bin/sh

set -eu

CONCOURSE_AWS_SSM_ACCESS_KEY=$(aws ssm get-parameters --name /concourse/main/sus/ci_aws_access_key_id --with-decryption --query "Parameters[0].Value" --output text) \
CONCOURSE_AWS_SSM_SECRET_KEY=$(aws ssm get-parameters --name /concourse/main/sus/ci_aws_secret_access_key --with-decryption --query "Parameters[0].Value" --output text) \
docker-compose up -d
