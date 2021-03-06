package s3

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3Service struct {
	s3iface.S3API
}

var GotPutObjectInput *s3.PutObjectInput

type MockObject struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Number  int    `json:"number"`
}

func (m mockS3Service) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	obj := &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(`
		{"name":"Proctor", "version":"2.0", "number":1}
		`))),
	}
	return obj, nil
}

func (m mockS3Service) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	GotPutObjectInput = input
	return &s3.PutObjectOutput{}, nil
}

func TestLoad(t *testing.T) {
	mc := Client{
		Service: mockS3Service{},
	}

	want := MockObject{
		Name:    "Proctor",
		Version: "2.0",
		Number:  1,
	}

	got := MockObject{}

	if err := mc.Load(&got, ""); err == nil {
		t.Fatal("Empty path - expected error, got: nil")
	}

	if err := mc.Load(&got, "s3/path"); err != nil {
		t.Fatalf("Expected error nil, got: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got: %#v; expected: %#v", got, want)
	}
}

func TestSave(t *testing.T) {
	mc := Client{
		Service: mockS3Service{},
		Bucket:  "Mockery",
	}

	mo := MockObject{
		Name:    "SaveMe",
		Version: "1.0",
		Number:  99,
	}

	if err := mc.Save(&mo, ""); err == nil {
		t.Fatal("Empty path - expected error, got: nil")
	}

	if err := mc.Save(&mo, "s3/path"); err != nil {
		t.Fatalf("Expected error nil, got: %v", err)
	}
	test, _ := json.MarshalIndent(mo, "", "  ")
	switch {
	case *GotPutObjectInput.Bucket != "Mockery":
		t.Fatal("Incorrect Bucket:", *GotPutObjectInput.Bucket)
	case *GotPutObjectInput.ContentType != "application/json":
		t.Fatal("Incorrect ContentType:", *GotPutObjectInput.ContentType)
	case *GotPutObjectInput.Key != "s3/path":
		t.Fatal("Incorrect Key:", *GotPutObjectInput.Key)
	case !reflect.DeepEqual(GotPutObjectInput.Body, bytes.NewReader(test)):
		t.Fatal("Incorrect Body")
	}
}
