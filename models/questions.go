package models

// Questions is a versioned collection of Question's
type Questions struct {
	List              map[string]Question `json:"questions"`
	Version           string              `json:"version"`
	RisklevelsVersion string              `json:"risklevels_version"`
	Updated           string              `json:"updated"`
}

// Question has information about a question
type Question struct {
	Text     string            `json:"text"`
	LongText string            `json:"long_text"`
	Answers  map[string]Answer `json:"answers"`
}

// Answer has information about an answer
type Answer struct {
	Text      string   `json:"text"`
	Datatypes []string `json:"datatypes"`
}

// Path returns the main S3 path containing questions versions for a campaign
func (ql Questions) Path(c string) string {
	return "questions/" + c + "/"
}

// Object returns the full S3 path to the object containing questions data for a specific campaign/version
func (ql Questions) Object(c, v string) string {
	if v == "" {
		return ""
	}
	return ql.Path(c) + v + "/questions.json"
}
