package tailfile

import (
	"fmt"
	"goLogCollection/common"
	"goLogCollection/kafka"
	"log"
	"time"

	"github.com/Shopify/sarama"

	"github.com/hpcloud/tail"
)

type tailTask struct {
	path    string
	topic   string
	tailObj *tail.Tail
}

// InitTail 初始化若干log文件配置项
func InitTail(allConf []common.CollectEntry) (err error) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}

	for _, conf := range allConf {
		tt := tailTask{
			path:  conf.Path,
			topic: conf.Topic,
		}
		tt.tailObj, err = tail.TailFile(conf.Path, config)
		if err != nil {
			log.Println("init tail failed ", err, conf.Path)
			continue
		}

		//启动后台goroutine收集日志
		go tt.run()
	}

	return nil
}

func (t *tailTask) run() (err error) {
	//tail->kafka
	log.Printf("colloct for path %s is running", t.path)
	for {
		line, ok := <-t.tailObj.Lines
		if !ok {
			log.Printf("tail file close reopen, path:%s\n", t.path)
			time.Sleep(time.Second) //读取出错等待一秒继续读
			continue
		}
		//空行不发送
		if len(line.Text) == 0 {
			continue
		}
		fmt.Println("msg : ", line.Text)

		//利用通道将同步代码改为异步,封装成kafka的msg信息发送
		msg := &sarama.ProducerMessage{}
		msg.Topic = t.topic
		msg.Value = sarama.StringEncoder(line.Text)
		kafka.RecvMsg(msg)
	}
}
