package models

// Responses is a collection of responses to questions submitted by the client
type Responses struct {
	List              map[string]string `json:"responses"`
	RisklevelsVersion string            `json:"risklevels_version"`
	QuestionsVersion  string            `json:"questions_version"`
	Metadata          interface{}       `json:"metadata"`
}

// Path returns the main S3 path containing responses for a campaign
func (r Responses) Path(c string) string {
	return "responses/" + c + "/"
}
