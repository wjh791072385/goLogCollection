package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

var Client sarama.SyncProducer

func InitKafka(address []string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	//连接kafka
	Client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Println("producer closed, err:", err)
		return err
	}

	return nil
}
