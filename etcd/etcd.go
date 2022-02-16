package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"goLogCollection/common"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var cli *clientv3.Client

func InitEtcd(address []string) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect to etcd failed, err: ", err)
		return err
	}
	return nil
}

func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := cli.Get(ctx, key)

	if err != nil {
		log.Println("get key failed ", err)
		return nil, err
	}

	//取第一个键值对
	if len(resp.Kvs) == 0 {
		log.Println("get key nil ", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Kvs[0].Value, &collectEntryList)
	return
}
