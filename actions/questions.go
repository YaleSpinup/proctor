package actions

import (
	"errors"

	"github.com/YaleSpinup/proctor/libs/helpers"
	"github.com/YaleSpinup/proctor/models"
	"github.com/gobuffalo/buffalo"
)

// QuestionsGet gets a list of questions for a given campaign
// Optional "version" query param can specify a version, otherwise the latest one will be used
func QuestionsGet(c buffalo.Context) error {
	questions := models.Questions{}
	campaign := c.Param("campaign")
	version := c.Param("version")

	// determine latest version if none specified
	if len(version) == 0 {
		vl, err := S3.GetVersions(questions.Path(campaign))
		if err != nil {
			return errors.New("Unable to determine latest questions version")
		}
		version, err = helpers.LatestVersion(vl)
		if err != nil {
			return errors.New("Unable to determine latest questions version")
		}
	}

	if err := S3.Load(&questions, questions.Object(campaign, version)); err != nil {
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
