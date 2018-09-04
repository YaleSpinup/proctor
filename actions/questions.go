package actions

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gobuffalo/buffalo"
)

// Answer has information about an answer
type Answer struct {
	Text      string   `json:"text"`
	Datatypes []string `json:"datatypes"`
}

// Question has information about a question
type Question struct {
	Text    string            `json:"text"`
	Answers map[string]Answer `json:"answers"`
}

// QuestionsList is a versioned collection of questions
type QuestionsList struct {
	List              map[string]Question `json:"questions"`
	Version           string              `json:"version"`
	RisklevelsVersion string              `json:"risklevels_version"`
	Updated           string              `json:"updated"`
}

// QuestionsGet gets a list of questions for a given campaign
// Optional "version" query param can specify a version, otherwise the latest one will be used
func QuestionsGet(c buffalo.Context) error {
	q, err := loadQuestions(c.Param("campaign"), c.Param("version"))
	if err != nil {
		if len(q) == 0 {
			return c.Error(404, err)
		}
		return err
	}

	var ql *QuestionsList
	if err := json.Unmarshal(q, &ql); err != nil {
		return c.Error(500, errors.New("Unable to unmarshal questions"))
	}

	// sanitize questions as to not show the data types in the answers
	for k, v := range ql.List {
		for x, y := range v.Answers {
			y.Datatypes = nil
			v.Answers[x] = y
		}
		ql.List[k] = v
	}

	return c.Render(200, r.JSON(ql))
}

// loadQuestions loads the questions json from S3 and returns a slice of bytes
func loadQuestions(campaign, version string) ([]byte, error) {
	if len(version) == 0 {
		// determine latest version
		vl, err := getVersions("questions/"+campaign+"/", "/")
		if err != nil {
			return nil, errors.New("Unable to determine latest questions version for " + campaign)
		}
		version = latestVersion(vl)
		if len(version) == 0 {
			return nil, errors.New("Unable to determine latest questions version for " + campaign)
		}
	}

	log.Println("Loading " + campaign + " questions version " + version)

	q, err := S3.GetObject("questions/" + campaign + "/" + version + "/questions.json")
	if err != nil {
		if len(q) == 0 {
			return []byte{}, errors.New("Object not found in S3")
		}
		return nil, errors.New("Unable to get object from S3")
	}

	return q, nil
}
