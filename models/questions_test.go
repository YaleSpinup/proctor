package models

import (
	"testing"
)

var questions = Questions{}

func TestPath(t *testing.T) {

	mystring := "abc123"
	got := questions.Path(mystring)
	badretval := "questions/abc123/"
	if got != badretval {
		t.Fatalf("Got unexpected string %s, expected %s", badretval, got)
	}

	/*
		//negative test
			retval := fmt.Sprintf("badquestions/%s/", mystring)
			if got != retval {
				//	t.Logf("Got unexpected string %s", retval)
				t.Errorf("Got expected string %s", retval)
			}
	*/
}

func TestObject(t *testing.T) {
	mystring := "foo"
	yourstring := "bar"

	got := questions.Object(yourstring, mystring)
	//fmt.Printf(got)
	if got == "mystring" {
		t.Errorf("Got expected string: %s", mystring)
	}
}
