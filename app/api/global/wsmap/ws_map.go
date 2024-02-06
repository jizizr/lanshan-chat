package wsmap

import (
	"container/list"
	"sync"
)

type WSMap struct {
	wsMap sync.Map
}

type Update struct {
	Action string
	Data   interface{}
}

func NewWSMap() *WSMap {
	return &WSMap{sync.Map{}}
}

func (w *WSMap) NewConnection(groupID int64, c chan *Update) *list.Element {
	v, _ := w.wsMap.LoadOrStore(groupID, list.New())
	return v.(*list.List).PushBack(c)
}

func (w *WSMap) GetConnections(groupID int64) *list.List {
	v, ok := w.wsMap.Load(groupID)
	if !ok {
		return nil
	}
	return v.(*list.List)
}
