package actions

import (
	"errors"
	"log"

	"github.com/YaleSpinup/proctor/libs/helpers"
	"github.com/YaleSpinup/proctor/models"
	"github.com/gobuffalo/buffalo"
)

// RiskLevelsGet gets a list of supported risk levels and data types
// Optional "version" query param can specify a version, otherwise the latest one will be used
func RiskLevelsGet(c buffalo.Context) error {
	risklevels := models.RiskLevels{}
	version := c.Param("version")

	// determine latest version if none specified
	if len(version) == 0 {
		vl, err := S3.GetVersions(risklevels.Path())
		if err != nil {
			return errors.New("Unable to determine latest risklevels version")
		}
		version, err = helpers.LatestVersion(vl)
		if err != nil {
			return errors.New("Unable to determine latest risklevels version")
		}
	}

	if err := S3.Load(&risklevels, risklevels.Object(version)); err != nil {
		return err
	}

	log.Println(risklevels)
	return c.Render(200, r.JSON(risklevels))
}
