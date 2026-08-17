#!/bin/sh
set -eu

#aws s3 mb s3://test-bucket
#aws sqs create-queue --queue-name testqueue
echo "List S3 buckets"
aws s3 ls
