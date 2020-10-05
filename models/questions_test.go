package models

import (
	"testing"
)

var questions = Questions{}

func TestPath(t *testing.T) {

	mystring := "abc123"
	got := questions.Path(mystring)
	retval := "questions/abc123/"
	if got != retval {
		t.Errorf("Got unexpected string %s, expected %s", retval, got)
	}
}

func TestObject(t *testing.T) {
	mystring := "foo"
	yourstring := "bar"

	got := questions.Object(yourstring, mystring)
	if got == "mystring" {
		t.Errorf("Got expected string: %s", mystring)
	}
}
