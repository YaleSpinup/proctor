package models

import (
	"errors"

	"github.com/YaleSpinup/proctor/libs/helpers"
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

// Path returns the main S3 path containing risk level versions
func (rl RiskLevels) Path() string {
	return "risklevels/"
}

// Object returns the full S3 path to the object containing risk level data for a specific version
func (rl RiskLevels) Object(v string) string {
	return rl.Path() + v + "/risklevels.json"
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
