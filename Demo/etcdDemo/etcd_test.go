package etcdDemo

import (
	"context"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func TestEtcdGetPut(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		t.Logf("connect to etcd failed, err:%v\n", err)
		return
	}
	t.Log("init success")
	defer cli.Close()

	//put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "etcdTes", "jike")
	if err != nil {
		t.Log("put error", err)
		return
	}
	cancel()

	//get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "etcdTes")
	cancel()
	if err != nil {
		t.Log("get error", err)
		return
	}

	//kv键值对遍历
	for _, ev := range resp.Kvs {
		t.Logf("%s 	%s", ev.Key, ev.Value)
	}
}

//watch操作监控etcd中key的变化
func TestWatch(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Logf("connect to etcd failed, err:%v\n", err)
		return
	}
	t.Log("init success")
	defer cli.Close()

	//watch监控key, 返回一个<-chan WatchResponse
	wch := cli.Watch(context.Background(), "etcdTes")
	for wresp := range wch {
		for _, event := range wresp.Events {
			t.Logf("Type: %s Key:%s Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
		}
	}

}
