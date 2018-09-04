package actions

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gobuffalo/buffalo"
)

// Response has the data types and risk level determined based on the questions response
type Response struct {
	DataTypes []string `json:"datatypes"`
	RiskLevel string   `json:"risklevel"`
}

// ResponsesList is a collection of responses to questions submitted by the client
type ResponsesList struct {
	List              map[string]string `json:"responses"`
	RisklevelsVersion string            `json:"risklevels_version"`
	QuestionsVersion  string            `json:"questions_version"`
}

// ResponsesPost processes the question responses and returns the data type and security level in a Response
func ResponsesPost(c buffalo.Context) error {
	rl := &ResponsesList{}
	if err := c.Bind(rl); err != nil {
		return c.Error(400, errors.New("Bad request"))
	}

	// get questions/answers for the given campaign and version
	var ql *QuestionsList
	q, err := loadQuestions(c.Param("campaign"), rl.QuestionsVersion)
	if err != nil {
		if len(q) == 0 {
			return c.Error(404, err)
		}
		return err
	}
	if err = json.Unmarshal(q, &ql); err != nil {
		return c.Error(500, errors.New("Unable to unmarshal questions"))
	}

	// determine the associated data types based on all the answers
	dt, err := processResponse(rl, ql)
	if err != nil {
		return c.Error(422, err)
	}

	// get mapping of data types to risk levels
	var risklevels *RiskLevelsList
	rls, err := loadRiskLevels(c.Param("version"))
	if err != nil {
		if len(rls) == 0 {
			return c.Error(404, err)
		}
		return err
	}
	if err := json.Unmarshal(rls, &risklevels); err != nil {
		return errors.New("Unable to unmarshal risk levels")
	}

	// build the response
	var resp Response
	resp.DataTypes = dt
	resp.RiskLevel = highestRiskLevel(dt, risklevels.List)

	return c.Render(200, r.JSON(resp))
}

// processResponse returns a slice of data types, or error if the response is invalid
func processResponse(rl *ResponsesList, ql *QuestionsList) ([]string, error) {
	// basic response length validation
	if len(rl.List) != len(ql.List) {
		return nil, errors.New("Invalid question response: all questions must be answered")
	}

	dt := []string{}
	for id, v := range ql.List {
		ra, ok := rl.List[id]
		if !ok {
			return nil, errors.New("Invalid question response: no answer for question " + id)
		}

		a, ok := v.Answers[ra]
		if !ok {
			return nil, errors.New("Invalid question response: invalid answer \"" + ra + "\" for question " + id)
		}
		dt = append(dt, a.Datatypes...)
	}

	return uniqueSlice(dt), nil
}

// highestRiskLevel determines the highest risk level based on a list of data types and a list of RiskLevels
func highestRiskLevel(dt []string, rl []RiskLevel) string {
	var score uint
	for _, d := range dt {
		for _, r := range rl {
			if stringInSlice(d, r.Datatypes) && r.Score > score {
				score = r.Score
			}
		}
	}
	for _, r := range rl {
		if r.Score == score {
			return r.Text
		}
	}
	log.Println("Error: Unable to determine risk level for data types", dt)
	return ""
}

// uniqueSlice returns a slice containing only the unique strings from the original
func uniqueSlice(s []string) []string {
	seen := map[string]bool{}
	for k := range s {
		seen[s[k]] = true
	}

	result := []string{}
	for k := range seen {
		result = append(result, k)
	}
	return result
}

// stringInSlice returns true if s is in the list slice
func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
