package tasks

import (
	lru "github.com/hashicorp/golang-lru"
	"meepo/util"
)

var ResultCache *lru.Cache

func init() {
	var err error
	ResultCache, err = lru.New(500)
	if err != nil {
		util.SugarLogger.Errorf("初始化缓存失败：%s", err)
	}
}

func GetResult(uuid string) []byte {
	var rst []byte
	if val, ok := ResultCache.Get(uuid); ok {
		switch v := val.(type) {
		case chan []byte:
			util.SugarLogger.Debugf("获取到chan []byte")
			rst = <-v
			ResultCache.Add(uuid, rst)
		case []byte:
			util.SugarLogger.Debugf("获取到[]byte")
			rst = v
		default:
			util.SugarLogger.Errorf("获取到UUID【%s】的结果不是【chan []byte】类型", uuid)
		}
	} else {
		util.SugarLogger.Errorf("获取不到UUID【%s】的结果", uuid)
	}
	return rst
}