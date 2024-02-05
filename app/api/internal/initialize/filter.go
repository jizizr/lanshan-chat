package initialize

import (
	"github.com/sgoware/go-sensitive"
	"go.uber.org/zap"
	"lanshan_chat/app/api/global"
)

func SetupFilter() {
	filterManager := sensitive.NewFilter(
		sensitive.StoreOption{
			Type: sensitive.StoreMemory,
		},
		sensitive.FilterOption{
			Type: sensitive.FilterDfa,
		},
	)

	// 加载字典
	err := filterManager.GetStore().LoadDictPath(global.Config.FilterConfig.DictPath)
	if err != nil {
		global.Logger.Fatal("load dict failed", zap.Error(err))
		return
	}
	global.Filter = filterManager
}
