package version

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	groupPrefix            string = "prefix"
	groupMajor             string = "major"
	groupMinor             string = "minor"
	groupPatch             string = "patch"
	groupRevision          string = "revision"
	groupPrerelease        string = "prerelease"
	groupPrereleasePrefix  string = "prerelease_prefix"
	groupPrereleaseCounter string = "prerelease_counter"
	groupMetadata          string = "build_metadata"
)

var expression = regexp.MustCompile(
	`^\s*(?:` +
		// Versions can start with an optional V
		fmt.Sprintf(`\s*(?<%s>v|V)?`, groupPrefix) +
		// Get the major version number
		fmt.Sprintf(`\s*(?<%s>\d+)\s*`, groupMajor) +
		// Get the minor version number, delimited by a period, comma, dash or underscore
		fmt.Sprintf(`(?:\.\s*(?<%s>\d+)\s*)?`, groupMinor) +
		// Get the patch version number, delimited by a period, comma, dash or underscore
		fmt.Sprintf(`(?:\.\s*(?<%s>\d+)\s*)?`, groupPatch) +
		// Get the revision version number, delimited by a period, comma, dash or underscore
		fmt.Sprintf(`(?:\.\s*(?<%s>\d+)\s*)?)?`, groupRevision) +
		// Everything after the last digit and before the plus is the prerelease
		fmt.Sprintf(
			`(?:[.\-_\\])?(?<%s>(?<%s>[A-Za-z0-9]*?)([.\-_\\](?<%s>[A-Za-z0-9.\-_\\]*?)?)?)?`,
			groupPrerelease,
			groupPrereleasePrefix,
			groupPrereleaseCounter,
		) +
		// The metadata is everything after the plus
		fmt.Sprintf(`(?:\+(?<%s>[A-Za-z0-9_\-.\\+]*?))?\s*$`, groupMetadata),
)

func ParseOctopusVersion(value string) (*OctopusVersion, error) {
	if len(value) == 0 {
		return nil, errors.New("the version cannot be an empty string")
	}

	noSpaces := strings.Replace(value, " ", "", -1)

	regexResult := expression.FindStringSubmatch(value)

	if regexResult == nil {
		return nil, errors.New("the supplied version was not valid")
	}
	mappedResult := make(map[string]string)
	for i, name := range expression.SubexpNames() {
		if i != 0 && name != "" {
			mappedResult[name] = regexResult[i]
		}
	}

	var major, minor, patch, revision int
	var err error
	if major, err = getInt32FromMapOrDefault(mappedResult, groupMajor, 0); err != nil {
		return nil, fmt.Errorf("failed to parse major version: %v", err)
	}
	if minor, err = getInt32FromMapOrDefault(mappedResult, groupMinor, 0); err != nil {
		return nil, fmt.Errorf("failed to parse minor version: %v", err)
	}
	if patch, err = getInt32FromMapOrDefault(mappedResult, groupPatch, 0); err != nil {
		return nil, fmt.Errorf("failed to parse patch version: %v", err)
	}
	if revision, err = getInt32FromMapOrDefault(mappedResult, groupRevision, 0); err != nil {
		return nil, fmt.Errorf("failed to parse revision: %v", err)
	}

	return NewOctopusVersion(
			getStringFromMapOrDefault(mappedResult, groupPrefix, ""),
			major,
			minor,
			patch,
			revision,
			getStringFromMapOrDefault(mappedResult, groupPrerelease, ""),
			getStringFromMapOrDefault(mappedResult, groupPrereleasePrefix, ""),
			getStringFromMapOrDefault(mappedResult, groupPrereleaseCounter, ""),
			getStringFromMapOrDefault(mappedResult, groupMetadata, ""),
			noSpaces),
		nil
}

func getStringFromMapOrDefault(m map[string]string, key string, defaultValue string) string {
	if val, ok := m[key]; ok {
		return val
	}
	return defaultValue
}

func getInt32FromMapOrDefault(m map[string]string, key string, defaultValue int) (int, error) {
	if val, ok := m[key]; ok && val != "" {
		// the C# only parses as an int32, so we do the same here.
		i, err := strconv.ParseInt(val, 10, 32)
		if err == nil {
			return int(i), nil
		}
		// we had an error
		return 0, err
	}
	return defaultValue, nil
}
