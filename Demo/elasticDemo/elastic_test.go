package elasticDemo

import (
	"context"
	"fmt"
	"testing"

	"github.com/olivere/elastic/v7"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func TestElastic(t *testing.T) {
	//如果ElasticSearch是通过docker安装，若不设置elastic.SetSniff(false)
	//会报错: no active connection found: no Elasticsearch node available
	//elastic.SetSniff(false) 设置是否定期检查集群（默认为true）
	cli, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"), elastic.SetSniff(false))
	if err != nil {
		t.Log("connect failed", err)
		return
	}

	t.Log("connect successfully")
	p1 := Person{Name: "lmh", Age: 18, Married: false}
	put1, err := cli.Index().Index("user").BodyJson(p1).Do(context.Background())
	if err != nil {
		t.Log("put failed")
		return
	}

	fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
