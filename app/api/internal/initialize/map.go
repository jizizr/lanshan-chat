package initialize

import (
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/global/wsmap"
)

func SetupMap() {
	global.WSMap = wsmap.NewWSMap()
}
