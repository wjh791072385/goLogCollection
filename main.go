package main

import (
	"goLogCollection/kafka"
	"goLogCollection/tailfile"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"

	"gopkg.in/ini.v1"
)

// Config global config
type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

func main() {
	//加载配置文件到结构体中
	cfg, err := ini.Load("./conf/config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	globalConfig := new(Config) //因为要传指针，所以采用new
	err = cfg.MapTo(globalConfig)
	if err != nil {
		log.Printf("Fail: %v", err)
		os.Exit(1)
	}

	err = kafka.InitKafka([]string{globalConfig.KafkaConfig.Address}, globalConfig.KafkaConfig.ChanSize)
	if err != nil {
		log.Println("fail init kafka", err)
		return
	}
	log.Println("kafka init success")

	err = tailfile.InitTail(globalConfig.CollectConfig.LogFilePath)
	if err != nil {
		log.Println("fail init tail", err)
		return
	}
	log.Println("tail init success")

	//将日志发往kafka
	run()

}

func run() (err error) {
	//tail->kafka
	for {
		line, ok := <-tailfile.TailObj.Lines
		if !ok {
			log.Printf("tail file close reopen, filename:%s\n", tailfile.TailObj.Filename)
			time.Sleep(time.Second) //读取出错等待一秒继续读
			continue
		}

		//利用通道将同步代码改为异步,封装成kafka的msg信息发送
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder(line.Text)
		kafka.MsgChan <- msg
	}
}
