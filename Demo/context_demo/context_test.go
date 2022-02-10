package context_demo

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

var wgc sync.WaitGroup

func workerc(ctx context.Context) {
	defer wgc.Done()
exitLable:
	for {
		fmt.Println("worker...")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break exitLable
		default:

		}
	}
}

func TestContextDemo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wgc.Add(1)
	go workerc(ctx)
	time.Sleep(3 * time.Second)
	cancel() //告诉子goroutine退出
	wgc.Wait()
	fmt.Println("over")
}
