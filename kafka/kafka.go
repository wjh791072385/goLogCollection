package kafka

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func InitKafka(address []string, chanSize int64) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	//连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Println("producer closed, err:", err)
		return err
	}

	//初始化msg通道
	msgChan = make(chan *sarama.ProducerMessage, chanSize)

	//持续读取消息,非main函数结束，不会导致goroutine退出
	go sendMsg()

	return nil
}

func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				log.Println("send msg fail ", err)
			}
			log.Println(pid, offset)
		}
	}
}

func RecvMsg(msg *sarama.ProducerMessage) {
	msgChan <- msg
	fmt.Println("send successfully")
}
