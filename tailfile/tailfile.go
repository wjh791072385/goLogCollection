package tailfile

import (
	"log"

	"github.com/hpcloud/tail"
)

var TailObj *tail.Tail

func InitTail(fileName string) (err error) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}

	TailObj, err = tail.TailFile(fileName, config)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
