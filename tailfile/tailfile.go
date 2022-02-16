package tailfile

import (
	"fmt"
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

func (t *tailTask) Init() (err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tailObj, err = tail.TailFile(t.path, cfg)
	return
}
