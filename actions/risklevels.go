package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gobuffalo/buffalo"
)

// RiskLevel describes a security risk level and what data types belong to it
type RiskLevel struct {
	Text      string   `json:"text"`
	Score     uint     `json:"score"`
	Datatypes []string `json:"datatypes"`
}

// RiskLevels is a versioned collection of RiskLevel's
type RiskLevels struct {
	List    []RiskLevel `json:"risklevels"`
	Version string      `json:"version"`
	Updated string      `json:"updated"`
}

// RiskLevelsGet gets a list of supported risk levels and data types
// Optional "version" query param can specify a version, otherwise the latest one will be used
func RiskLevelsGet(c buffalo.Context) error {
	rl, err := loadRiskLevels(c.Param("version"))
	if err != nil {
		if len(rl) == 0 {
			return c.Error(404, err)
		}
		return err
	}

	var risklevels RiskLevels
	if err := json.Unmarshal(rl, &risklevels); err != nil {
		return errors.New("Unable to unmarshal risk levels")
	}

	log.Println(risklevels)
	return c.Render(200, r.JSON(risklevels))
}

// loadRisklevels loads the risk levels json from S3 and returns a slice of bytes
func loadRiskLevels(version string) ([]byte, error) {
	if len(version) == 0 {
		// determine latest version
		vl, err := getVersions("risklevels/", "/")
		if err != nil {
			return nil, errors.New("Unable to determine latest risklevels version")
		}
		version = latestVersion(vl)
		if len(version) == 0 {
			return nil, errors.New("Unable to determine latest risklevels version")
		}
	}

	log.Printf("Loading risk levels version %s", version)
	path := fmt.Sprintf("risklevels/%s/risklevels.json", version)
	d, err := S3.GetObject(path)
	if err != nil {
		if len(d) == 0 {
			return []byte{}, errors.New("Object not found in S3")
		}
		return nil, errors.New("Unable to get object from S3")
	}

	return d, nil
}
