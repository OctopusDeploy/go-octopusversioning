package version

type ComparableVersion[T any] interface {
	Compare(T) int
}

func Compare[T ComparableVersion[T]](v1, v2 T) int {
	return v1.Compare(v2)
}
