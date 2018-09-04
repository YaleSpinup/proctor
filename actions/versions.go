package actions

import (
	"sort"
	"strings"
)

// getVersions returns a list of all versions under a given prefix
// This assumes a specific format when organizing the objects in the S3 bucket, e.g.
//   s3://bucket/prefix/1.0/file.json
func getVersions(prefix, delimiter string) ([]string, error) {
	dv, err := S3.ListObjects(prefix, delimiter)
	if err != nil {
		return nil, err
	}

	var versions []string
	// split the object prefixes to get the last part which is the version
	for _, p := range dv.CommonPrefixes {
		s := strings.Split(strings.Trim(*p.Prefix, "/"), "/")
		versions = append(versions, s[len(s)-1])
	}
	return versions, nil
}

// latestVersion gets a slice of versions, e.g. [0.1 0.2 1.0 1.1] and returns the latest version as a string
func latestVersion(vl []string) string {
	if len(vl) == 0 {
		return ""
	}
	sort.Strings(vl)
	return vl[len(vl)-1]
}
