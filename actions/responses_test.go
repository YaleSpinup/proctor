package actions

import (
	"encoding/json"
	"reflect"

	"github.com/YaleSpinup/proctor/models"
)

func (as *ActionSuite) Test_ResponsesPost() {
	o := Outcome{}
	r := models.Responses{}

	// empty response
	got := as.JSON("/v1/proctor/test/responses").Post(r)
	as.Equal(400, got.Code)

	// invalid response (missing List)
	r = models.Responses{
		List: map[string]string{"1": "a"},
	}
	got = as.JSON("/v1/proctor/test/responses").Post(r)
	as.Equal(400, got.Code)

	// invalid response (missing RisklevelsVersion)
	r = models.Responses{
		List:             map[string]string{"1": "a"},
		QuestionsVersion: "2.0",
	}
	got = as.JSON("/v1/proctor/test/responses").Post(r)
	as.Equal(400, got.Code)

	// invalid response (missing QuestionsVersion)
	r = models.Responses{
		List:              map[string]string{"1": "a"},
		RisklevelsVersion: "1.0",
	}
	got = as.JSON("/v1/proctor/test/responses").Post(r)
	as.Equal(400, got.Code)

	// invalid response (unknown question)
	r = models.Responses{
		List:              map[string]string{"99": "a"},
		RisklevelsVersion: "1.0",
		QuestionsVersion:  "2.0",
	}
	got = as.JSON("/v1/proctor/test/responses").Post(r)
	as.Equal(422, got.Code)

	// invalid response (unknown answer)
	r = models.Responses{
		List:              map[string]string{"1": "x"},
		RisklevelsVersion: "1.0",
		QuestionsVersion:  "2.0",
	}
	got = as.JSON("/v1/proctor/test/responses").Post(r)
	as.Equal(422, got.Code)

	// valid response - low risk
	wantDataTypes := []string{}
	wantRiskLevel := "low"
	r = models.Responses{
		List:              map[string]string{"1": "b"},
		RisklevelsVersion: "1.0",
		QuestionsVersion:  "2.0",
	}

	got = as.JSON("/v1/proctor/test/responses").Post(r)
	as.Equal(200, got.Code)
	err := json.Unmarshal([]byte(got.Body.String()), &o)
	if err != nil {
		as.T().Fatalf("Got response: %#v; Error: %v", got.Body.String(), err)
	}
	if !reflect.DeepEqual(o.DataTypes, wantDataTypes) {
		as.T().Fatalf("Got: %#v; expected: %#v", o.DataTypes, wantDataTypes)
	}
	if o.RiskLevel != wantRiskLevel {
		as.T().Fatalf("Got: %#v; expected: %#v", o.RiskLevel, wantRiskLevel)
	}

	// valid response - high risk
	wantDataTypes = []string{"HIPAA"}
	wantRiskLevel = "high"
	r = models.Responses{
		List:              map[string]string{"1": "a"},
		RisklevelsVersion: "1.0",
		QuestionsVersion:  "2.0",
		Metadata:          "blah",
	}

	got = as.JSON("/v1/proctor/test/responses").Post(r)
	as.Equal(200, got.Code)
	err = json.Unmarshal([]byte(got.Body.String()), &o)
	if err != nil {
		as.T().Fatalf("Got response: %#v; Error: %v", got.Body.String(), err)
	}
	if !reflect.DeepEqual(o.DataTypes, wantDataTypes) {
		as.T().Fatalf("Got: %#v; expected: %#v", o.DataTypes, wantDataTypes)
	}
	if o.RiskLevel != wantRiskLevel {
		as.T().Fatalf("Got: %#v; expected: %#v", o.RiskLevel, wantRiskLevel)
	}
}
