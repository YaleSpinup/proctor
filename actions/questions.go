package actions

import (
	"github.com/YaleSpinup/proctor/models"
	"github.com/gobuffalo/buffalo"
)

// QuestionsGet gets a list of questions for a given campaign
// Optional "version" query param can specify a version, otherwise the latest one will be used
func QuestionsGet(c buffalo.Context) error {
	questions := models.Questions{}
	if err := questions.Load(&S3, c.Param("campaign"), c.Param("version")); err != nil {
		return err
	}

	// sanitize questions as to not show the data types in the answers
	for k, v := range questions.List {
		for x, y := range v.Answers {
			y.Datatypes = nil
			v.Answers[x] = y
		}
		questions.List[k] = v
	}

	return c.Render(200, r.JSON(questions))
}
