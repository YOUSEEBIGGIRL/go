package mytest

import (
	"fmt"
	"testing"
	"time"
)

func TestChanBuf(t *testing.T) {
	ch := make(chan int64, 3)

	ch <- 1
	ch <- 2
	ch <- 3

	go func() {
		for v := range ch {
			fmt.Println(v)
		}
	}()

	time.Sleep(time.Second)
}

func TestChan(t *testing.T) {

}

func TestNilChan(t *testing.T) {
	var ch chan int64
	ch <- 1	// 向 nil chan 写入数据，会导致当前 goroutine 被阻塞

	go func() {
		<-ch	// 可以从 nil chan 中读取数据
	}()
}
