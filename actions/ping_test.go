package actions

func (as *ActionSuite) Test_PingPong() {
	res := as.JSON("/v1/proctor/ping").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "pong")
}
