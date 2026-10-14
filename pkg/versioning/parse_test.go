package versioning

import (
	"testing"
)

func TestOctopusVersionParsing(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		major          int
		minor          int
		patch          int
		revision       int
		release        string
		releasePrefix  string
		releaseCounter string
		metadata       string
	}{
		{"0.0.4", "0.0.4", 0, 0, 4, 0, "", "", "", ""},
		{"0.0.4.1", "0.0.4.1", 0, 0, 4, 1, "", "", "", ""},
		{"1.2.3", "1.2.3", 1, 2, 3, 0, "", "", "", ""},
		{"1.2.3.4", "1.2.3.4", 1, 2, 3, 4, "", "", "", ""},
		{"10.20.30", "10.20.30", 10, 20, 30, 0, "", "", "", ""},
		{"v0.0.0", "v0.0.0", 0, 0, 0, 0, "", "", "", ""},
		{"v0.0.1", "v0.0.1", 0, 0, 1, 0, "", "", "", ""},
		{"v1.0.0", "v1.0.0", 1, 0, 0, 0, "", "", "", ""},
		{"v0.0.0-foo", "v0.0.0-foo", 0, 0, 0, 0, "foo", "foo", "", ""},
		{"1.1.2-prerelease+meta", "1.1.2-prerelease+meta", 1, 1, 2, 0, "prerelease", "prerelease", "", "meta"},
		{"1.1.2+meta", "1.1.2+meta", 1, 1, 2, 0, "", "", "", "meta"},
		{"1.1.2+meta-valid", "1.1.2+meta-valid", 1, 1, 2, 0, "", "", "", "meta-valid"},
		{"1.0.0-alpha", "1.0.0-alpha", 1, 0, 0, 0, "alpha", "alpha", "", ""},
		{"1.0.0-alpha.beta", "1.0.0-alpha.beta", 1, 0, 0, 0, "alpha.beta", "alpha", "beta", ""},
		{"1.0.0-alpha.beta1", "1.0.0-alpha.beta1", 1, 0, 0, 0, "alpha.beta1", "alpha", "beta1", ""},
		{"1.0.0-alpha.1", "1.0.0-alpha.1", 1, 0, 0, 0, "alpha.1", "alpha", "1", ""},
		{"1.0.0-alpha0.valid", "1.0.0-alpha0.valid", 1, 0, 0, 0, "alpha0.valid", "alpha0", "valid", ""},
		{"1.0.0-alpha.0valid", "1.0.0-alpha.0valid", 1, 0, 0, 0, "alpha.0valid", "alpha", "0valid", ""},
		{
			"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
			"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
			1,
			0,
			0,
			0,
			"alpha-a.b-c-somethinglong",
			"alpha",
			"a.b-c-somethinglong",
			"build.1-aef.1-its-okay",
		},
		{"1.0.0-rc.1+build.1", "1.0.0-rc.1+build.1", 1, 0, 0, 0, "rc.1", "rc", "1", "build.1"},
		{"2.0.0-rc.1+build.123", "2.0.0-rc.1+build.123", 2, 0, 0, 0, "rc.1", "rc", "1", "build.123"},
		{"1.2.3-beta", "1.2.3-beta", 1, 2, 3, 0, "beta", "beta", "", ""},
		{"10.2.3-DEV-SNAPSHOT", "10.2.3-DEV-SNAPSHOT", 10, 2, 3, 0, "DEV-SNAPSHOT", "DEV", "SNAPSHOT", ""},
		{"1.2.3-SNAPSHOT-123", "1.2.3-SNAPSHOT-123", 1, 2, 3, 0, "SNAPSHOT-123", "SNAPSHOT", "123", ""},
		{"1.0.0", "1.0.0", 1, 0, 0, 0, "", "", "", ""},
		{"2.0.0", "2.0.0", 2, 0, 0, 0, "", "", "", ""},
		{"1.1.7", "1.1.7", 1, 1, 7, 0, "", "", "", ""},
		{"2.0.0+build.1848", "2.0.0+build.1848", 2, 0, 0, 0, "", "", "", "build.1848"},
		{"2.0.1-alpha.1227", "2.0.1-alpha.1227", 2, 0, 1, 0, "alpha.1227", "alpha", "1227", ""},
		{"1.0.0-alpha+beta", "1.0.0-alpha+beta", 1, 0, 0, 0, "alpha", "alpha", "", "beta"},
		{
			"1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
			"1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
			1,
			2,
			3,
			0,
			"---RC-SNAPSHOT.12.9.1--.12",
			"",
			"--RC-SNAPSHOT.12.9.1--.12",
			"788",
		},
		{
			"1.2.3----R-S.12.9.1--.12+meta",
			"1.2.3----R-S.12.9.1--.12+meta",
			1,
			2,
			3,
			0,
			"---R-S.12.9.1--.12",
			"",
			"--R-S.12.9.1--.12",
			"meta",
		},
		{
			"1.2.3----RC-SNAPSHOT.12.9.1--.12",
			"1.2.3----RC-SNAPSHOT.12.9.1--.12",
			1,
			2,
			3,
			0,
			"---RC-SNAPSHOT.12.9.1--.12",
			"",
			"--RC-SNAPSHOT.12.9.1--.12",
			"",
		},
		{
			"1.0.0+0.build.1-rc.10000aaa-kk-0.1",
			"1.0.0+0.build.1-rc.10000aaa-kk-0.1",
			1,
			0,
			0,
			0,
			"",
			"",
			"",
			"0.build.1-rc.10000aaa-kk-0.1",
		},
		{"99999999.99999999.99999999", "99999999.99999999.99999999", 99999999, 99999999, 99999999, 0, "", "", "", ""},
		{"0.0.0-foo", "0.0.0-foo", 0, 0, 0, 0, "foo", "foo", "", ""},
		{"0.0.0", "0.0.0", 0, 0, 0, 0, "", "", "", ""},
		{"1.2", "1.2", 1, 2, 0, 0, "", "", "", ""},
		{"1", "1", 1, 0, 0, 0, "", "", "", ""},
		{"01.1.1", "01.1.1", 1, 1, 1, 0, "", "", "", ""},
		{"1.01.1", "1.01.1", 1, 1, 1, 0, "", "", "", ""},
		{"1.1.01", "1.1.01", 1, 1, 1, 0, "", "", "", ""},
		{
			"1.1.01-9999999999999999999999999",
			"1.1.01-9999999999999999999999999",
			1,
			1,
			1,
			0,
			"9999999999999999999999999",
			"9999999999999999999999999",
			"",
			"",
		},
		{
			"trailing whitespace after metadata",
			"1.1.01-9999999999999999999999999+meta ",
			1,
			1,
			1,
			0,
			"9999999999999999999999999",
			"9999999999999999999999999",
			"",
			"meta",
		},
		{
			"trailing whitespace",
			"1.1.01-9999999999999999999999999 ",
			1,
			1,
			1,
			0,
			"9999999999999999999999999",
			"9999999999999999999999999",
			"",
			"",
		},
		{
			"leading whitespace",
			" 1.1.01-9999999999999999999999999 ",
			1,
			1,
			1,
			0,
			"9999999999999999999999999",
			"9999999999999999999999999",
			"",
			"",
		},
		{
			"whitespace around major",
			" 1 .1.01-9999999999999999999999999",
			1,
			1,
			1,
			0,
			"9999999999999999999999999",
			"9999999999999999999999999",
			"",
			"",
		},
		{
			"whitespace around minor",
			"1. 1 .01-9999999999999999999999999",
			1,
			1,
			1,
			0,
			"9999999999999999999999999",
			"9999999999999999999999999",
			"",
			"",
		},
		{
			"whitespace around patch",
			"1.1. 01 -9999999999999999999999999",
			1,
			1,
			1,
			0,
			"9999999999999999999999999",
			"9999999999999999999999999",
			"",
			"",
		},
		{
			"whitespace around release",
			"1.1.01. 1 -9999999999999999999999999",
			1,
			1,
			1,
			1,
			"9999999999999999999999999",
			"9999999999999999999999999",
			"",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseOctopusVersion(tt.version)
			if err != nil {
				t.Fatalf("ParseOctopusVersion(%q) failed: %v", tt.version, err)
			}

			if parsed.Major != tt.major {
				t.Errorf("Major = %d, want %d", parsed.Major, tt.major)
			}
			if parsed.Minor != tt.minor {
				t.Errorf("Minor = %d, want %d", parsed.Minor, tt.minor)
			}
			if parsed.Patch != tt.patch {
				t.Errorf("Patch = %d, want %d", parsed.Patch, tt.patch)
			}
			if parsed.Revision != tt.revision {
				t.Errorf("Revision = %d, want %d", parsed.Revision, tt.revision)
			}
			if parsed.Release != tt.release {
				t.Errorf("Release = %q, want %q", parsed.Release, tt.release)
			}
			if parsed.ReleasePrefix != tt.releasePrefix {
				t.Errorf("ReleasePrefix = %q, want %q", parsed.ReleasePrefix, tt.releasePrefix)
			}
			if parsed.ReleaseCounter != tt.releaseCounter {
				t.Errorf("ReleaseCounter = %q, want %q", parsed.ReleaseCounter, tt.releaseCounter)
			}
			if parsed.Metadata != tt.metadata {
				t.Errorf("Metadata = %q, want %q", parsed.Metadata, tt.metadata)
			}
		})
	}
}
