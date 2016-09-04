package s3proxy

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sheeley/s3-object-proxy/permission"
)

type ImageProxy struct {
	bucket           string
	s3               *s3.S3
	permissionLookup permission.Lookup
}

func New(c *Config) (*ImageProxy, error) {
	if err := c.Setup(); err != nil {
		return nil, err
	}
	return &ImageProxy{
		s3:               c.S3,
		bucket:           c.Bucket,
		permissionLookup: c.PermissionLookup,
	}, nil
}

func (i *ImageProxy) GetObject(key string) ([]byte, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(key),
	}

	if !i.canServe(key) {
		return nil, errors.New("404 file not found")
	}

	log.Println("Getting object from S3")
	resp, err := i.s3.GetObject(params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Object received")
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (i *ImageProxy) canServe(key string) bool {
	response, err := i.permissionLookup.CanView(key)
	if err == nil && response.Public {
		return true
	}
	if err != nil {
		log.Println(err)
	}
	return false
}
