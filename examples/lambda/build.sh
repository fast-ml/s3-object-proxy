#! /usr/bin/env bash

if [[ "$BUCKET" == "" ]]; then
  echo "BUCKET is required"
  exit 1
fi

if [[ "$REGION" == "" ]]; then
  echo "REGION is required"
  exit 1
fi

if [[ "$REDIS_HOST" == "" ]]; then
  echo "REDIS_HOST is required"
  exit 1
fi

OUTPUT_FILE="proxy"
ZIP_FILE="lambda-proxy.zip"
LD_FLAGS="-X main.bucket=$BUCKET -X main.region=$REGION -X main.redisAddr=$REDIS_HOST"
set -ex 

rm -f ./$ZIP_FILE
GOOS=linux GOARCH=amd64 go build -ldflags "$LD_FLAGS" -o $OUTPUT_FILE
zip -r $ZIP_FILE $OUTPUT_FILE index.js 
