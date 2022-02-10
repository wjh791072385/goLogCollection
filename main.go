package main

import (
	"goLogCollection/kafka"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// Config global config
type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
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

	err = kafka.InitKafka([]string{globalConfig.KafkaConfig.Address})
	if err != nil {
		log.Println("fail init kafka", err)
	}

}
