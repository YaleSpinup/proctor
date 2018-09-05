package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/YaleSpinup/proctor/libs/helpers"
	"github.com/YaleSpinup/proctor/libs/s3"
)

// Questions is a versioned collection of Question's
type Questions struct {
	List              map[string]Question `json:"questions"`
	Version           string              `json:"version"`
	RisklevelsVersion string              `json:"risklevels_version"`
	Updated           string              `json:"updated"`
}

// Question has information about a question
type Question struct {
	Text    string            `json:"text"`
	Answers map[string]Answer `json:"answers"`
}

// Answer has information about an answer
type Answer struct {
	Text      string   `json:"text"`
	Datatypes []string `json:"datatypes"`
}

// Load loads the risk levels json from S3 and returns a slice of bytes
func (q *Questions) Load(s3 *s3.Client, campaign, version string) error {
	if len(version) == 0 {
		// determine latest version
		vl, err := s3.GetVersions("questions/"+campaign+"/", "/")
		if err != nil {
			return errors.New("Unable to determine latest questions version for " + campaign)
		}
		version = helpers.LatestVersion(vl)
		if len(version) == 0 {
			return errors.New("Unable to determine latest questions version for " + campaign)
		}
	}

	log.Printf("Loading %s questions version %s", campaign, version)
	path := fmt.Sprintf("questions/%s/%s/questions.json", campaign, version)
	o, err := s3.GetObject(path)
	if err != nil {
		if len(o) == 0 {
			return errors.New("Object not found in S3")
		}
		return errors.New("Unable to get object from S3")
	}

	if err := json.Unmarshal(o, q); err != nil {
		return errors.New("Unable to unmarshal questions")
	}

	return nil
}
