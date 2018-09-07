package models

// Responses is a collection of responses to questions submitted by the client
type Responses struct {
	List              map[string]string `json:"responses"`
	RisklevelsVersion string            `json:"risklevels_version"`
	QuestionsVersion  string            `json:"questions_version"`
}
