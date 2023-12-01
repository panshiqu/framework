package utils

import (
	"math/rand"
)

type Builtin interface {
	Integer | ~string
}

type Integer interface {
	~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64
}

func RandSlice[T Builtin](s []T) T {
	index := rand.Intn(len(s))
	return s[index]
}

func RandBetween[T Integer](a, b T) T {
	if a >= b {
		return a
	}
	// for safe use
	if T(int32(a)) != a || T(int32(b)) != b {
		return a
	}
	return a + T(rand.Int63n(int64(b-a)))
}

func RandWeight[T Integer](s []T) int {
	var sum int64
	areas := make([]int64, len(s))
	for k, v := range s {
		// for safe use
		if T(uint16(v)) != v {
			return -1
		}
		sum += int64(v)
		areas[k] = sum
	}
	if sum == 0 {
		return -1
	}
	n := rand.Int63n(sum)
	for k, v := range areas {
		if n < v {
			return k
		}
	}
	return -1
}
