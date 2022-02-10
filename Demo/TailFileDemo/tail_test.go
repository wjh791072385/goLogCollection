package TailFileDemo

import (
	"testing"
	"time"

	"github.com/hpcloud/tail"
)

func TestTail(t *testing.T) {
	fileName := "./my.log"

	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}

	tf, err := tail.TailFile(fileName, config)
	if err != nil {
		t.Log(err)
		return
	}

	//开始读取数据
	var line *tail.Line
	var ok bool

	for {
		line, ok = <-tf.Lines
		if !ok {
			t.Logf("tail file close reopen, filename:%s\n", tf.Filename)
			time.Sleep(time.Second) //读取出错
			continue
		}
		t.Log("line:", line.Text, line.Time)
	}
}
