package service

import (
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/global/wsmap"
)

func Notify(groupID int64, update *wsmap.Update) {
	list := global.WSMap.GetConnections(groupID)
	if list == nil {
		return
	}
	for e := list.Front(); e != nil; e = e.Next() {
		ch := e.Value.(chan *wsmap.Update)
		go func() { ch <- update }()
	}
}
