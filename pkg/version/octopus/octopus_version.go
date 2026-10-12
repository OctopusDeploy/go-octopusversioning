package version

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"

	"github.com/OctopusDeploy/go-octopusversioning/pkg/version"
	utils "github.com/OctopusDeploy/go-octopusversioning/pkg/version/util"
)

type OctopusVersion struct {
	Prefix         string
	Major          int
	Minor          int
	Patch          int
	Revision       int
	Release        string
	ReleasePrefix  string
	ReleaseCounter string
	Metadata       string
	OriginalString string
}

func NewOctopusVersion(
	prefix string,
	major int,
	minor int,
	patch int,
	revision int,
	release string,
	releasePrefix string,
	releaseCounter string,
	metadata string,
	originalString string,
) *OctopusVersion {
	return &OctopusVersion{
		Prefix:         prefix,
		Major:          major,
		Minor:          minor,
		Patch:          patch,
		Revision:       revision,
		Release:        release,
		ReleasePrefix:  releasePrefix,
		ReleaseCounter: releaseCounter,
		Metadata:       metadata,
		OriginalString: originalString,
	}
}

func (o *OctopusVersion) Compare(other version.ComparableVersion) (int, error) {
	switch otherTyped := other.(type) {
	case *OctopusVersion:
		return o.CompareTo(otherTyped), nil
	default:
		return 0, fmt.Errorf("cannot compare OctopusVersion to %s", otherTyped)
	}
}

func (o *OctopusVersion) CompareTo(other *OctopusVersion) int {
	if majorCmp := cmp.Compare(o.Major, other.Major); majorCmp != 0 {
		return majorCmp
	}
	if minorCmp := cmp.Compare(o.Minor, other.Minor); minorCmp != 0 {
		return minorCmp
	}
	if patchCmp := cmp.Compare(o.Patch, other.Patch); patchCmp != 0 {
		return patchCmp
	}
	if revisionCmp := cmp.Compare(o.Revision, other.Revision); revisionCmp != 0 {
		return revisionCmp
	}

	// anything with a release field is lower than anything without a release field
	if len(o.Release) > 0 && len(other.Release) == 0 {
		return -1
	}
	if len(other.Release) > 0 && len(o.Release) == 0 {
		return 1
	}

	releaseAlphaNum := utils.GetAlphaNumericOnly(o.Release)
	otherReleaseAlphaNum := utils.GetAlphaNumericOnly(other.Release)

	// if the releases are different
	if releaseAlphaNum != otherReleaseAlphaNum {
		return compareReleaseLabels(
			strings.FieldsFunc(releaseAlphaNum, splitter),
			strings.FieldsFunc(otherReleaseAlphaNum, splitter))
	}

	return 0
}

func compareReleaseLabels(v1 []string, v2 []string) int {
	maxLen := max(len(v1), len(v2))
	for i := range maxLen {
		// v2 is longer
		if i >= len(v1) && i < len(v2) {
			return -1
		}
		// v1 is longer
		if i < len(v1) && i >= len(v2) {
			return 1
		}

		result := compareRelease(v1[i], v2[i])

		if result != 0 {
			return result
		}
	}
	return 0
}

func compareRelease(v1 string, v2 string) int {
	i1, err1 := strconv.Atoi(v1)
	i2, err2 := strconv.Atoi(v2)

	// if both numeric, compare as numeric
	if err1 == nil && err2 == nil {
		return cmp.Compare(i1, i2)
	} else if err1 == nil || err2 == nil {
		// numeric labels come before alpha labels
		if err1 == nil {
			return -1
		}

		return 1
	}

	// Ignoring 2.0.0 case sensitive compare. Everything will be compared case insensitively as 2.0.1 specifies.
	return strings.Compare(strings.ToLower(v1), strings.ToLower(v2))
}

func splitter(r rune) bool {
	return r == '-' || r == '_' || r == '.'
}
