package utils

import (
	"cmp"
	"slices"
)

func Pointer[T any](value T) *T {
	return &value
}

func SlicesEqualSorted[T cmp.Ordered](a, b []T) bool {
	slices.Sort(a)
	slices.Sort(b)
	return slices.Equal(a, b)
}
