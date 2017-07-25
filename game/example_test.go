package game

import (
	"math/rand"
	"sort"
	"testing"
)

func TestTableFrameSlice(t *testing.T) {
	var m []*TableFrame
	m = append(m, &TableFrame{id: 1, count: 3})
	m = append(m, &TableFrame{id: 2, count: 3})
	m = append(m, &TableFrame{id: 3, count: 2})
	m = append(m, &TableFrame{id: 4, count: 2})
	m = append(m, &TableFrame{id: 5, count: 1})
	m = append(m, &TableFrame{id: 6, count: 1})
	m = append(m, &TableFrame{id: 7})
	m = append(m, &TableFrame{id: 8})
	m = append(m, &TableFrame{id: 9, count: 4})
	m = append(m, &TableFrame{id: 10, count: 4})
	m = append(m, &TableFrame{id: 11, status: 1})
	m = append(m, &TableFrame{id: 12, status: 1})

	mm := make([]*TableFrame, len(m))

	for n := 0; n < 100000; n++ {
		perm := rand.Perm(len(m))

		for k, v := range perm {
			mm[k] = m[v]
		}

		sort.Sort(TableFrameSlice(mm))

		for i := 0; i < len(m); i++ {
			if m[i] != mm[i] {
				t.Error(m[i], mm[i])
			}
		}
	}
}

func TestTrySitDown(t *testing.T) {

}
