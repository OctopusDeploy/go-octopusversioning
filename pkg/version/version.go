package version

type ComparableVersion interface {
	Compare(ComparableVersion) (int, error)
}

func Compare[T ComparableVersion](v1, v2 T) (int, error) {
	return v1.Compare(v2)
}
