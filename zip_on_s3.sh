# build
make build

# zip without dependecy third files
zip -r /tmp/sus.zip ./ \
    -x *.git* vendor/\* integration_test/vendor/\* ci/params.yml

# upload to s3 bucket
aws s3 cp /tmp/sus.zip s3://$S3_BUCKET

# clean
make clean
