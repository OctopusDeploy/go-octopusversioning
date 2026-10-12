package version

import "regexp"

var nonAlphaNumeric = regexp.MustCompile("[^\\p{L}\\p{N}]")

func GetAlphaNumericOnly(s string) string {
	return nonAlphaNumeric.ReplaceAllString(s, "-")
}
