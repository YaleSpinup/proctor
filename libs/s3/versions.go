package s3

import "strings"

// GetVersions returns a list of all versions under a given prefix in S3
// This assumes a specific format when organizing the objects in the S3 bucket, e.g.
//   s3://bucket/prefix/1.0/file.json
func (s Client) GetVersions(prefix string) ([]string, error) {
	dv, err := s.ListObjects(prefix, "/")
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
