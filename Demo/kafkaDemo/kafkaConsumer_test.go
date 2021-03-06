package kafkaDemo

import (
	"fmt"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

//消费者
func TestKafkaConsumer(t *testing.T) {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}

	partitionList, err := consumer.Partitions("web_log") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	t.Log(partitionList)

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetOldest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		//defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				t.Logf("Partition:%d Offset:%d Key:%s Value:%s", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}

	//阻塞方便测试
	time.Sleep(100 * time.Second)
}
