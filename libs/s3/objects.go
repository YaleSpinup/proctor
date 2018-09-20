package s3

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

// Load loads the json object from S3 (path) into the given interface
func (s Client) Load(i interface{}, path string) error {
	if len(path) == 0 {
		return errors.New("Path cannot be empty")
	}

	log.Printf("Loading %s", path)
	o, err := s.GetObject(path)
	if err != nil {
		if len(o) == 0 {
			return errors.New("Object not found in S3")
		}
		return errors.New("Unable to get object from S3")
	}

	if err := json.Unmarshal(o, i); err != nil {
		return fmt.Errorf("Unable to unmarshal %T", i)
	}

	return nil
}

// PutObject takes a byte array and saves it in S3 with the given key name
func (s Client) PutObject(object []byte, key string, ct string) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(object),
		ContentType: aws.String(ct),
	}

	if _, err := s.Service.PutObject(input); err != nil {
		log.Println(err.Error(), input)
		return err
	}

	return nil
}

// Save saves an object to S3 (path) as JSON
func (s Client) Save(i interface{}, path string) error {
	if len(path) == 0 {
		return errors.New("Path cannot be empty")
	}

	log.Printf("Saving %T to %s", i, path)
	j, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Errorf("Unable to marshal %T", i)
	}

	if err = s.PutObject(j, path, "application/json"); err != nil {
		return errors.New("Unable to put object into S3")
	}

	return nil
}
