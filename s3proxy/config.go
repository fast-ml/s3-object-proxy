package s3proxy

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sheeley/s3-object-proxy/permission"
)

// Config contains AWS configuration and PermissionLookups
// Bucket, PermissionLookups, and one of (S3, ASConfig, or Region) are required
type Config struct {
	Bucket           string
	Region           string
	AWSConfig        *aws.Config
	Session          *session.Session
	S3               *s3.S3
	PermissionLookup permission.Lookup
}

func appendIf(test bool, errs []string, err string) []string {
	if test {
		return append(errs, err)
	}
	return errs
}

func (c *Config) Validate() error {
	errs := []string{}
	appendIf(c.Bucket == "", errs, "Bucket is required")
	appendIf(c.PermissionLookup == nil, errs, "PermissionLookup is required")
	appendIf(c.S3 == nil && c.AWSConfig == nil && c.Region == "", errs, "One of S3, AWSConfig, or Region is required")

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func (c *Config) Setup() error {
	if err := c.Validate(); err != nil {
		return err
	}

	if c.S3 == nil {
		if c.Session == nil {
			c.Session = session.New()
		}

		if c.AWSConfig == nil {
			if c.Region == "" {
				return errors.New("One of S3, AWSConfig or region is required")
			}
			c.AWSConfig = &aws.Config{Region: aws.String(c.Region)}
		}

		c.S3 = s3.New(c.Session, c.AWSConfig)
	}
	return nil
}
