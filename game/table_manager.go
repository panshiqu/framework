package game

import (
	"sort"
	"sync"
)

var tins TableManager

// TableManager 桌子管理
type TableManager struct {
	count  int
	mutex  sync.Mutex
	tables []*TableFrame
}

// TrySitDown 尝试坐下
func (t *TableManager) TrySitDown(userItem *UserItem) {
	for {
		sort.Sort(TableFrameSlice(t.tables))

		// 只要有桌子椅子就能坐下，这里不关心桌子状态
		if len(t.tables) == 0 || t.tables[0].UserCount() == cins.UserPerTable {
			t.AddTableFrame()
			continue
		}

		t.tables[0].SitDown(userItem)
		break
	}
}

// AddTableFrame 增加桌子
func (t *TableManager) AddTableFrame() {
	t.count++

	tableFrame := &TableFrame{
		id: t.count,
	}

	t.tables = append(t.tables, tableFrame)
}

// TableFrameSlice 排序
type TableFrameSlice []*TableFrame

func (t TableFrameSlice) Len() int {
	return len(t)
}
func (t TableFrameSlice) Less(i, j int) bool {
	if t[i].TableStatus() != t[j].TableStatus() {
		return t[i].TableStatus() < t[j].TableStatus()
	} else if c1, c2 := t[i].UserCount(), t[j].UserCount(); c1 != c2 {
		switch {
		case c2 == cins.UserPerTable:
			return true
		case c1 == cins.UserPerTable:
			return false
		default:
			return c1 > c2
		}
	}

	return t[i].TableID() < t[j].TableID()
}
func (t TableFrameSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
