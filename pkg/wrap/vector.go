package wrap

import (
	"golang.org/x/exp/constraints"
)

func VectorIndex[T constraints.Ordered](list []T, v T) int {
	for i := range list {
		if list[i] == v {
			return i
		}
	}

	return -1
}

func EqualVector[T constraints.Ordered](arr1, arr2 []T) (ok bool) {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func UniqVector[T constraints.Ordered](arr []T) (list []T) {
	n := len(arr)
	list = make([]T, 0, n)

	if len(arr) == 0 {
		return list
	}

	mp := make(map[T]bool, n)
	for _, v := range arr {
		if !mp[v] {
			list = append(list, v)
			mp[v] = true
		}
	}

	return list
}

func First[T any](v []T) *T {
	if len(v) == 0 {
		return nil
	}
	return &v[0]
}

func Last[T any](v []T) *T {
	var n = len(v)

	if n == 0 {
		return nil
	}
	return &v[n-1]
}
