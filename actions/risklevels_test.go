package actions

func (as *ActionSuite) Test_RiskLevelsGet() {
	wantBody := `{"risklevels":[{"text":"high","score":30,"datatypes":["HIPAA","PCI","SSN"]},{"text":"moderate","score":20,"datatypes":["FERPA"]},{"text":"low","score":0,"datatypes":[]}],"version":"1.0","updated":""}`
	wantCode := 200

	got := as.JSON("/v1/proctor/risklevels").Get()
	as.Equal(wantCode, got.Code)
	as.Contains(got.Body.String(), wantBody)
}
