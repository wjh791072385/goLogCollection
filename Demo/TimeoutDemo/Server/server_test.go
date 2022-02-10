package Server

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func serverHandler(w http.ResponseWriter, r *http.Request) {
	number := rand.Intn(2)
	if number == 1 {
		time.Sleep(time.Second * 10) // 耗时10秒的慢响应
		fmt.Fprintf(w, "slow response")
		return
	}
	fmt.Fprint(w, "quick response")
}

func TestServer(t *testing.T) {
	http.HandleFunc("/server", serverHandler)

	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		panic(err)
	}
}
