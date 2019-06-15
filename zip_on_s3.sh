# prevent uploading big binary file
make clean

# zip without dependecy third files
zip -r sus.zip ./ \
  -x *.git* vendor/\* integration_test/vendor/\* ci/params.yml

# upload to s3 bucket
aws s3 cp sus.zip s3://$S3_BUCKET

# remove zipped file
rm sus.zip
