package actions

import (
	"log"

	"github.com/YaleSpinup/proctor/models"
	"github.com/gobuffalo/buffalo"
)

// RiskLevelsGet gets a list of supported risk levels and data types
// Optional "version" query param can specify a version, otherwise the latest one will be used
func RiskLevelsGet(c buffalo.Context) error {
	risklevels := models.RiskLevels{}
	if err := risklevels.Load(&S3, c.Param("version")); err != nil {
		return err
	}

	log.Println(risklevels)
	return c.Render(200, r.JSON(risklevels))
}
