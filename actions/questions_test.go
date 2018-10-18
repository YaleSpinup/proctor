package actions

func (as *ActionSuite) Test_QuestionsGet() {
	wantBody := `{"questions":{"1":{"text":"Do you have HIPAA data?","long_text":"","answers":{"a":{"text":"Yes","datatypes":null},"b":{"text":"No","datatypes":null}}}},"version":"2.0","risklevels_version":"1.0","updated":""}`
	wantCode := 200

	got := as.JSON("/v1/proctor/test/questions").Get()
	as.Equal(wantCode, got.Code)
	as.Contains(got.Body.String(), wantBody)
}
