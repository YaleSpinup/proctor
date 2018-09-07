package helpers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// LatestVersion gets a slice of versions and returns the latest version as a string
// e.g. ["0.1", "0.10", "1.0"] or ["0.1.0", "0.10.1", "1.0.2"]
// any number of minor/patch versions will work as long as all versions are consistent
// i.e. this is not allowed ["0.1", "0.10.0", "1.0"] and the function will return error
func LatestVersion(vl []string) (string, error) {
	if len(vl) == 0 {
		return "", errors.New("Unable to determine latest version, empty slice")
	}

	// initialize latest with zeros
	latest := []int{}
	for i := 0; i < len(strings.Split(vl[0], ".")); i++ {
		latest = append(latest, 0)
	}

	// loop over list and compare each to latest
	for _, sv := range vl {
		current := strings.Split(sv, ".")
		if len(current) < 2 {
			log.Printf("Error: format mismatch for version '%s'", sv)
			return "", errors.New("Unable to determine latest version, format mismatch")
		}

		si := []int{}
		for _, s := range current {
			i, err := strconv.Atoi(s)
			if err != nil {
				log.Printf("Error: format mismatch for version '%s'", sv)
				return "", errors.New("Unable to determine latest version, format mismatch")
			}
			si = append(si, i)
		}

		if len(latest) != len(si) {
			log.Printf("Error: format mismatch for version '%s'", sv)
			return "", errors.New("Unable to determine latest version, format mismatch")
		}

		for i, siv := range si {
			// we break out unless the version index is equal to the latest one
			if siv > latest[i] {
				latest = si
				break
			}
			if siv < latest[i] {
				break
			}
		}
	}
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(latest)), "."), "[]"), nil
}

// StringInSlice returns true if s is in the list slice
func StringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// UniqueSlice returns a slice containing only the unique strings from the original
func UniqueSlice(s []string) []string {
	seen := map[string]bool{}
	for k := range s {
		seen[s[k]] = true
	}

	result := []string{}
	for k := range seen {
		result = append(result, k)
	}
	return result
}
