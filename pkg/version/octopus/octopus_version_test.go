package version

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusversioning/pkg/version"
)

type VersionTest struct {
	v1             string
	v2             string
	expectedResult int
}

func TestVersionComparisons(t *testing.T) {
	tests := []VersionTest{
		{"1.0.0", "2.0.0", -1},
		{"1.1.0", "1.2.0", -1},
		{"1.1.1", "1.1.2", -1},
		{"1.1.1.1", "1.1.1.2", -1},
		{"1.1.0", "2.2.0", -1},
		{"1.1.1", "2.1.2", -1},
		{"1.1.1.1", "2.1.1.2", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.2.0", "1.1.0", 1},
		{"1.1.2", "1.1.1", 1},
		{"1.1.1.2", "1.1.1.1", 1},
		{"2.2.0", "1.1.0", 1},
		{"2.1.2", "1.1.1", 1},
		{"2.1.1.2", "1.1.1.1", 1},
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.1", "1.0.2", -1},
		{"1.1.0", "1.1.0", 0},
		{"1.1.1", "1.1.0", 1},
		{"1.1.1", "1.1.2", -1},
		{"1.0.0.0", "1.0.0.0", 0},
		{"1.0.1.0", "1.0.0.0", 1},
		{"1.0.1.0", "1.0.2.0", -1},
		{"1.0.0-prerelease", "1.0.0", -1},
		{"1.0.0", "1.0.0-prerelease", 1},
		{"1.0.0", "1.0.0-prerelease", 1},
		{"1.0.0-prerelease", "1.0.0-prerelease", 0},
		{"1.0.0-prerelease.1", "1.0.0-prerelease.1", 0},
		{"1.0.0-prerelease-1", "1.0.0-prerelease.1", 0},
		{"1.0.0-prerelease_1", "1.0.0-prerelease.1", 0},
		{"1.0.0-prerelease_1", "1.0.0-prerelease.2", -1},
		{"1.0.0-prerelease_3", "1.0.0-prerelease.2", 1},
		{"1.0.0-prerelease1", "1.0.0-prerelease2", -1},
		{"1.0.0-prerelease3", "1.0.0-prerelease2", 1},
		{"1.0.0+meta", "1.0.0+meta", 0},
		{"1.0.0+meta1", "1.0.0+meta2", 0},
		{"1.0.0+meta3", "1.0.0+meta2", 0},
		{"1.0.0.1", "1.0.0-2", 1},
	}

	runTests(t, tests)
}

func TestPrereleaseRules(t *testing.T) {
	tests := []VersionTest{
		{"1.0.0-test.1", "1.0.0-test.2", -1},
		{"1.0.0-alpha", "1.0.0-alpha.1", -1},
		{"1.0.0-alpha.1", "1.0.0-alpha.beta", -1},
		{"1.0.0-alpha.beta", "1.0.0-beta", -1},
		{"1.0.0-beta", "1.0.0-beta.2", -1},
		{"1.0.0-beta.2", "1.0.0-beta.11", -1},
		{"1.0.0-beta.11", "1.0.0-rc.1", -1},
		{"1.0.0-rc.1", "1.0.0", -1},
		{"1.0.0-test.1", "1.0.0-test-2", -1},
		{"1.0.0-test.1", "1.0.0-test_2", -1},
		{"1.0.0-test-1", "1.0.0-test-2", -1},
		{"1.0.0-test-1", "1.0.0-test_2", -1},
		{"1.0.0-test_1", "1.0.0-test-2", -1},
		{"1.0.0-test_1", "1.0.0-test_2", -1},
		{"1.0.0-alpha", "1.0.0-alpha_1", -1},
		{"1.0.0-beta.2", "1.0.0-beta-11", -1},
		{"1.0.0-beta.2", "1.0.0-beta_11", -1},
		{"1.0.0-beta-2", "1.0.0-beta-11", -1},
		{"1.0.0-beta-2", "1.0.0-beta_11", -1},
		{"1.0.0-beta_2", "1.0.0-beta-11", -1},
		{"1.0.0-beta_2", "1.0.0-beta_11", -1},
	}

	runTests(t, tests)
}

func TestVersionsWithEquivalentChars(t *testing.T) {
	tests := []VersionTest{
		{`1.1.1-prerelease_-.\`, `1.1.1-prerelease\.-_`, 0},
		{`1.1.1-\_.10`, `1.1.1._\.11`, -1},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []VersionTest) {
	for _, tt := range tests {
		symbol := "="
		switch tt.expectedResult {
		case 1:
			symbol = ">"
			break
		case -1:
			symbol = "<"
			break
		}

		t.Run(fmt.Sprintf("%s %s %s", tt.v1, symbol, tt.v2), func(t *testing.T) {
			parsedV1, err := ParseOctopusVersion(tt.v1)
			if err != nil {
				t.Fatalf("ParseOctopusVersion(%q) failed: %v", tt.v1, err)
			}

			parsedV2, err := ParseOctopusVersion(tt.v2)
			if err != nil {
				t.Fatalf("ParseOctopusVersion(%q) failed: %v", tt.v2, err)
			}

			result := version.Compare(parsedV1, parsedV2)

			if result != tt.expectedResult {
				t.Errorf("CompareTo(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expectedResult)
			}
		})
	}
}
