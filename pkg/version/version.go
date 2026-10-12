package version

type Format int

const (
	Semver        Format = iota
	Maven         Format = iota
	Docker        Format = iota
	Octopus       Format = iota
	Lexicographic Format = iota
)

type ComparableVersion interface {
	CompareTo(ComparableVersion) int
}

func Compare[T ComparableVersion](v1, v2 T) int {
	return v1.CompareTo(v2)
}
