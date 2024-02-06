package router

import (
	"container/list"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/global/wsmap"
	"lanshan_chat/app/api/internal/controller"
	"lanshan_chat/app/api/internal/dao/mysql"
	"net/http"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(c *gin.Context) {
	uid, ok := controller.GetUID(c)
	if !ok {
		return
	}
	groups, err := mysql.QueryAllGroup(uid)
	if err != nil {
		return
	}
	if len(groups) == 0 {
		return
	}
	wsConn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	updates := make(chan *wsmap.Update, 100)
	elements := make([]*list.Element, len(groups), len(groups))
	for i, groupID := range groups {
		elements[i] = global.WSMap.NewConnection(groupID, updates)
	}
	go func() {
		for {
			update := <-updates
			go func() {
				err := wsConn.WriteJSON(update)
				if err != nil {
					closeChan(groups, elements, updates)
					_ = wsConn.Close()
					return
				}
			}()
		}
	}()
}

func closeChan(
	groups []int64,
	elements []*list.Element,
	updates chan *wsmap.Update,
) {
	for i, groupID := range groups {
		global.WSMap.GetConnections(groupID).Remove(elements[i])
	}
	for {
		select {
		case <-updates:
		default:
			close(updates)
			return
		}
	}
}
