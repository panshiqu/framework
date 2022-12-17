package game

import (
	"log"
	"net/http"
	"sync"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/game/fiveinarow"
	"github.com/panshiqu/framework/game/landlords"
)

var tins TableManager

// TableManager 桌子管理
type TableManager struct {
	count  int           // 计数
	mutex  sync.Mutex    // 加锁
	tables []*TableFrame // 桌子
}

// GetTable 获取桌子
func (t *TableManager) GetTable(id int) *TableFrame {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if id < len(t.tables) {
		return t.tables[id]
	}
	return nil
}

// TrySitDown 尝试坐下
func (t *TableManager) TrySitDown(userItem *UserItem) (tableFrame *TableFrame) {
	var max int
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, v := range t.tables {
		if v.TableStatus() != define.TableStatusFree {
			continue
		}

		userCount := v.UserCount()
		if userCount == define.CG.UserPerTable {
			continue
		}

		if tableFrame == nil || userCount > max {
			max, tableFrame = userCount, v
		}

		if userCount+1 == define.CG.UserPerTable {
			break
		}
	}

	if tableFrame == nil {
		tableFrame = t.AddTableFrame()
	}

	tableFrame.SitDown(userItem)

	return
}

// AddTableFrame 增加桌子
func (t *TableManager) AddTableFrame() (tableFrame *TableFrame) {
	tableFrame = &TableFrame{
		id:    t.count,
		users: make([]*UserItem, define.CG.UserPerTable),
	}

	tableFrame.SetTableLogic(CreateTableLogic(tableFrame))

	t.tables = append(t.tables, tableFrame)

	t.count++

	return
}

// Monitor 监视器
func (t *TableManager) Monitor(w http.ResponseWriter, r *http.Request) {
	t.mutex.Lock()
	for _, v := range t.tables {
		v.Monitor(w, r)
	}
	t.mutex.Unlock()
}

// CreateTableLogic 创建桌子逻辑
func CreateTableLogic(v define.ITableFrame) (ret define.ITableLogic) {
	defer func() {
		if ret == nil {
			log.Fatal("CreateTableLogic fatal")
		}
	}()

	switch define.CG.GameType {
	case define.GameLandlords: // 斗地主
		return landlords.NewTableLogic(v)
	case define.GameFiveInARow: // 五子棋
		return fiveinarow.NewTableLogic(v)
	}

	return nil
}
