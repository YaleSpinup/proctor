package actions

import (
	"errors"
	"log"

	"github.com/YaleSpinup/proctor/libs/helpers"
	"github.com/YaleSpinup/proctor/models"
	"github.com/gobuffalo/buffalo"
)

// ResponsesPost processes the question responses and returns the data type and security level in a Response
func ResponsesPost(c buffalo.Context) error {
	responses := models.Responses{}
	if err := c.Bind(&responses); err != nil {
		return c.Error(400, errors.New("Bad request"))
	}
	// TODO: Add responses validation

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

	// build the response
	var resp models.Response
	resp.DataTypes = datatypes
	resp.RiskLevel = hr.Text
	log.Println("Response outcome", resp)

	return c.Render(200, r.JSON(resp))
}
