package context_demo

import (
	"fmt"
	"sync"
	"testing"
	"time"
)


var wg sync.WaitGroup


func worker(ch <-chan struct{}) {
	defer wg.Done()
	exitLable:
	for {
		select {
			case <-ch :
				break exitLable
		default:
			fmt.Println("worker...")
			time.Sleep(time.Second)
		}
	}


}

func TestDemo(t *testing.T) {
	exitChan := make(chan struct{})
	wg.Add(1)
	go worker(exitChan)
	time.Sleep(3 * time.Second)
	exitChan <- struct{}{}
	wg.Wait()
	fmt.Println("over")
}
