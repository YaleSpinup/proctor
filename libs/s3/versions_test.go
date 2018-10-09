package s3

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3Service struct {
	s3iface.S3API
}

func (m mockS3Service) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	prefix := input.Prefix
	obj := &s3.ListObjectsOutput{
		CommonPrefixes: []*s3.CommonPrefix{
			{Prefix: aws.String(*prefix + "1.0/")},
			{Prefix: aws.String(*prefix + "1.1/")},
			{Prefix: aws.String(*prefix + "2.0/")},
			{Prefix: aws.String(*prefix + "2.1/")},
			{Prefix: aws.String(*prefix + "2.2/")},
			{Prefix: aws.String(*prefix + "2.3/")},
			{Prefix: aws.String(*prefix + "3.0/")},
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

func TestGetVersions(t *testing.T) {
	mc := Client{
		Service: mockS3Service{},
	}

	want := []string{"1.0", "1.1", "2.0", "2.1", "2.2", "2.3", "3.0"}
	got, err := mc.GetVersions("questions/test/")
	if err != nil {
		t.Fatalf("Expected error nil, got: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got: %v; expected: %v", got, want)
	}
}

func BenchmarkGetVersions(b *testing.B) {
	mc := Client{
		Service: mockS3Service{},
	}

	for i := 0; i < b.N; i++ {
		mc.GetVersions("questions/test/")
	}
}
