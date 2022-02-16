package tailfile

import (
	"goLogCollection/common"
	"log"
)

type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask
	CollectEntryList []common.CollectEntry
	confChan         chan []common.CollectEntry
}

var (
	ttMgr *tailTaskMgr
)

// InitTail 初始化若干log文件配置项
func InitTail(allConf []common.CollectEntry) (err error) {
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		CollectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry),
	}

	for _, conf := range allConf {
		tt := tailTask{
			path:  conf.Path,
			topic: conf.Topic,
		}

		err = tt.Init()
		if err != nil {
			log.Println("tailTask init failed ", err)
			continue
		}
		//添加到tailTaskMap中
		ttMgr.tailTaskMap[tt.path] = &tt

		//启动后台goroutine收集日志
		go tt.run()
	}

	go ttMgr.watch() //后台监听配置
	return nil
}

func (t *tailTaskMgr) watch() {
	newConf := <-t.confChan
	log.Println("get new conf from etcd ：", newConf)

	//对新配置的增删改查
	for _, conf := range newConf {
		_, ok := t.tailTaskMap[conf.Path]
		if ok {
			continue
		}

		//不存在则新建tailTask
		tt := tailTask{
			path:  conf.Path,
			topic: conf.Topic,
		}

		err := tt.Init()
		if err != nil {
			log.Println("tailTask init failed ", err)
			continue
		}
		//添加到tailTaskMap中
		ttMgr.tailTaskMap[tt.path] = &tt

		//启动后台goroutine收集日志
		go tt.run()
	}
}

func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
