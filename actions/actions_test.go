package actions

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	proctorS3 "github.com/YaleSpinup/proctor/libs/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite"
)

type ActionSuite struct {
	*suite.Action
}

type mockS3Service struct {
	s3iface.S3API
}

func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), packr.New("../fixtures", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	// override the S3 client
	S3 = proctorS3.Client{
		Service: mockS3Service{},
		Bucket:  "Mockery",
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}

func (m mockS3Service) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	var s3GetObjectOutput string

	prefix := strings.Split(*input.Key, "/")[0]
	switch prefix {
	case "questions":
		s3GetObjectOutput = `
		{
			"questions": {
				"1": {
					"text": "Do you have HIPAA data?",
					"answers": {
						"a": { "text": "Yes", "datatypes": ["HIPAA"] },
						"b": { "text": "No", "datatypes": [] }
					}
				}
			},
			"version": "2.0",
			"risklevels_version": "1.0"
		}
		`
	case "risklevels":
		s3GetObjectOutput = `
		{
			"version": "1.0",
			"risklevels": [
				{ "score": 30, "text": "high", "datatypes": ["HIPAA", "PCI", "SSN"] },
				{ "score": 20, "text": "moderate", "datatypes": ["FERPA"] },
				{ "score": 0, "text": "low", "datatypes": [] }
			]
		}
		`
	}
	obj := &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(s3GetObjectOutput))),
	}
	return obj, nil
}

func (m mockS3Service) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	prefix := input.Prefix
	obj := &s3.ListObjectsOutput{
		CommonPrefixes: []*s3.CommonPrefix{
			{Prefix: aws.String(*prefix + "1.0/")},
			{Prefix: aws.String(*prefix + "1.1/")},
		},
		Contents: []*s3.Object{
			{Key: prefix},
		},
		Prefix:      prefix,
		Delimiter:   aws.String("/"),
		IsTruncated: aws.Bool(false),
	}
	return obj, nil
}

func (m mockS3Service) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return &s3.PutObjectOutput{}, nil
}
