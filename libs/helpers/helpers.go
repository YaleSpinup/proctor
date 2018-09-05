package helpers

import (
	"log"
	"strconv"
	"strings"
)

// LatestVersion gets a slice of versions and returns the latest version as a string
// e.g. ["0.1", "0.10", "1.0"] or ["0.1.0", "0.10.1", "1.0.2"]
// any number of minor/patch versions will work as long as all versions are consistent
// i.e. this is not allowed ["0.1", "0.10.0", "1.0"] and the function will return ""
func LatestVersion(vl []string) string {
	if len(vl) == 0 {
		return ""
	}

	var ss [][]string
	var vlen int

	// split all versions into their atomic parts (e.g. "1.0" becomes ["1", "0"]) and put into slice of slices ss
	for _, v := range vl {
		sv := strings.Split(v, ".")
		if vlen == 0 {
			vlen = len(sv)
		} else {
			if vlen != len(sv) {
				log.Printf("Error: format mismatch for version '%s'", v)
				return ""
			}
		}
		ss = append(ss, sv)
	}

	// bubble sort the slices starting with the minor version (the one on the right) and moving to the left
	for l := vlen - 1; l >= 0; l-- {
		for m := 0; m < len(ss); m++ {
			for n := 0; n < len(ss)-1; n++ {
				int1, err := strconv.Atoi(ss[n][l])
				if err != nil {
					log.Printf("Error: format mismatch for version '%s'", strings.Join(ss[n], "."))
					return ""
				}
				int2, err := strconv.Atoi(ss[n+1][l])
				if err != nil {
					log.Printf("Error: format mismatch for version '%s'", strings.Join(ss[n+1], "."))
					return ""
				}
				if int1 > int2 {
					ss[n], ss[n+1] = ss[n+1], ss[n]
				}
			}
		}
	}

	return strings.Join(ss[len(ss)-1], ".")
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
