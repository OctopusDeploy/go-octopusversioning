package version

type Comparator[A, B any] func(A, B) int

func Compare[A, B any](v1 A, v2 B, compare Comparator[A, B]) int {
	return compare(v1, v2)
}
