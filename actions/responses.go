package actions

import (
	"errors"
	"log"

	"github.com/YaleSpinup/proctor/libs/helpers"
	"github.com/YaleSpinup/proctor/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/uuid"
)

// Outcome has the data types and risk level outcome determined based on the questions response
type Outcome struct {
	ID        uuid.UUID `json:"id"`
	DataTypes []string  `json:"datatypes"`
	RiskLevel string    `json:"risklevel"`
}

// ResponsesPost processes the question responses and returns the data type and security level in an Outcome
func ResponsesPost(c buffalo.Context) error {
	responses := models.Responses{}
	if err := c.Bind(&responses); err != nil {
		return c.Error(400, errors.New("Bad request"))
	}
	if len(responses.List) == 0 || len(responses.QuestionsVersion) == 0 || len(responses.RisklevelsVersion) == 0 {
		return c.Error(400, errors.New("Bad request"))
	}

	// get questions/answers for the given campaign and version
	questions := models.Questions{}
	if err := S3.Load(&questions, questions.Object(c.Param("campaign"), responses.QuestionsVersion)); err != nil {
		return err
	}

	// basic response length validation
	if len(responses.List) != len(questions.List) {
		return c.Error(422, errors.New("Invalid question response: all questions must be answered"))
	}

	// determine the associated data types based on all the answers in the response
	datatypes := []string{}
	for id, v := range questions.List {
		ra, ok := responses.List[id]
		if !ok {
			return c.Error(422, errors.New("Invalid question response: no answer for question "+id))
		}

		a, ok := v.Answers[ra]
		if !ok {
			return c.Error(422, errors.New("Invalid question response: invalid answer '"+ra+"' for question "+id))
		}
		datatypes = append(datatypes, a.Datatypes...)
	}
	datatypes = helpers.UniqueSlice(datatypes)

	// get mapping of data types to risk levels
	risklevels := models.RiskLevels{}
	if err := S3.Load(&risklevels, risklevels.Object(responses.RisklevelsVersion)); err != nil {
		return err
	}
	hr, err := risklevels.Highest(datatypes)
	if err != nil {
		return err
	}

	var outcome Outcome
	outcome.DataTypes = datatypes
	outcome.RiskLevel = hr.Text
	outcome.ID, err = uuid.NewV4()
	if err != nil {
		return c.Error(500, errors.New("Error generating outcome id"))
	}
	log.Println("Outcome response:", outcome)

	// save responses to S3
	// we also include the original questions and the outcome that was returned to the client
	s := struct {
		Outcome   Outcome
		Responses models.Responses
		Questions models.Questions
	}{outcome, responses, questions}
	path := responses.Path(c.Param("campaign")) + outcome.ID.String() + ".json"
	if err := S3.Save(s, path); err != nil {
		return err
	}

	return c.Render(200, r.JSON(outcome))
}
