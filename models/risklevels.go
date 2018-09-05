package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/YaleSpinup/proctor/libs/helpers"
	"github.com/YaleSpinup/proctor/libs/s3"
)

// RiskLevels is a versioned collection of RiskLevel's
type RiskLevels struct {
	List    []RiskLevel `json:"risklevels"`
	Version string      `json:"version"`
	Updated string      `json:"updated"`
}

// RiskLevel describes a security risk level and what data types belong to it
type RiskLevel struct {
	Text      string   `json:"text"`
	Score     uint     `json:"score"`
	Datatypes []string `json:"datatypes"`
}

// Load loads the risk levels json from S3 and returns a slice of bytes
func (rl *RiskLevels) Load(s3 *s3.Client, version string) error {
	if len(version) == 0 {
		// determine latest version
		vl, err := s3.GetVersions("risklevels/", "/")
		if err != nil {
			return errors.New("Unable to determine latest risklevels version")
		}
		version = helpers.LatestVersion(vl)
		if len(version) == 0 {
			return errors.New("Unable to determine latest risklevels version")
		}
	}

	log.Printf("Loading risk levels version %s", version)
	path := fmt.Sprintf("risklevels/%s/risklevels.json", version)
	o, err := s3.GetObject(path)
	if err != nil {
		if len(o) == 0 {
			return errors.New("Object not found in S3")
		}
		return errors.New("Unable to get object from S3")
	}

	if err := json.Unmarshal(o, rl); err != nil {
		return errors.New("Unable to unmarshal risk levels")
	}

	return nil
}

// Highest returns the RiskLevel with the highest score based on a list of data types
func (rl RiskLevels) Highest(datatypes []string) (RiskLevel, error) {
	var score uint
	for _, d := range datatypes {
		for _, r := range rl.List {
			if helpers.StringInSlice(d, r.Datatypes) && r.Score > score {
				score = r.Score
			}
		}
	}
	for _, r := range rl.List {
		if r.Score == score {
			return r, nil
		}
	}
	return RiskLevel{}, errors.New("Error: Unable to determine risk level for data types")
}
