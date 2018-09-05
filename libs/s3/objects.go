package s3

import (
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetObject returns an object from S3
func (s Client) GetObject(key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	result, err := s.Service.GetObject(input)
	if err != nil {
		log.Println(err.Error(), input)
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			return []byte{}, err
		}
		return nil, err
	}

	defer result.Body.Close()

	b, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return b, nil
}

// ListObjects returns a list of objects from S3
func (s Client) ListObjects(prefix, delimiter string) (*s3.ListObjectsOutput, error) {
	input := &s3.ListObjectsInput{
		Bucket:    aws.String(s.Bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delimiter),
	}

	result, err := s.Service.ListObjects(input)
	if err != nil {
		log.Println(err.Error(), input)
		return nil, err
	}

	return result, nil
}
