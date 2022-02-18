package main

import (
	"goLogCollection/etcd"
	"goLogCollection/kafka"
	"goLogCollection/tailfile"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// Config global config
type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `ini:"chan_size"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
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

	//传指针用new
	globalConfig := new(Config)

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

	//初始化etcd, 从etcd中拉取日志收集配置项
	err = etcd.InitEtcd([]string{globalConfig.EtcdConfig.Address})
	if err != nil {
		log.Println("fail init etcd", err)
		return
	}
	log.Println("etcd init success")

	allFilePath, err := etcd.GetLogConf(globalConfig.EtcdConfig.CollectKey)
	if err != nil {
		log.Println("fail get etcd_log_conf", err)
		return
	}

	//etcd监控配置文件
	go etcd.WatchLogConf(globalConfig.EtcdConfig.CollectKey)

	//读取配置文件对应log
	err = tailfile.InitTail(allFilePath)
	if err != nil {
		log.Println("fail init tail", err)
		return
	}

	select {}

}
