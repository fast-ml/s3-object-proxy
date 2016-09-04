package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sheeley/s3-object-proxy/permission"
	"github.com/sheeley/s3-object-proxy/s3proxy"
)

var bucket string
var region string
var redisAddr string

func main() {
	b, r, addr := getConfig()
	proxy, err := s3proxy.New(&s3proxy.Config{
		Bucket:           b,
		Region:           r,
		PermissionLookup: permission.NewRedisLookup(addr),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(proxy.ListenAndServe(":8080"))
}

func getConfig() (string, string, string) {
	flag.ErrHelp = errors.New("Usage: -bucket=bucket.name -region=bucket.region -redis=127.0.0.1:6379")
	_bucket := flag.String("bucket", bucket, "S3 bucket name")
	_region := flag.String("region", region, "S3 bucket region")
	_redisAddr := flag.String("redis", redisAddr, "Redis host:port")
	flag.Parse()

	if *_bucket == "" || *_region == "" || *_redisAddr == "" {
		fmt.Println(flag.ErrHelp)
		os.Exit(1)
	}

	return *_bucket, *_region, *_redisAddr
}
