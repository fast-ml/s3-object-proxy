package main

import "github.com/sheeley/s3-object-proxy/lambda"

var bucket string
var region string
var redisAddr string

func main() {
	lambda.DefaultRedisHandleFunc(bucket, region, redisAddr)
}
