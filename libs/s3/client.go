package s3

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Client struct contains the initialized S3 service and other S3-related parameters
type S3Client struct {
	Service *s3.S3
	Bucket  string
}

// NewSession creates an AWS session for S3 and returns an S3Client
func NewSession(key, secret, region string) S3Client {
	s := S3Client{}
	log.Printf("Creating new session with key id %s in region %s", key, region)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Region:      aws.String(region),
	}))
	s.Service = s3.New(sess)

	return s
}
