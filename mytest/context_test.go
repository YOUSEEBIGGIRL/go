package mytest

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContextCancel(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*5, cancelFunc)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("cancel!")
				return
			default:
				fmt.Println("do word...")
				time.Sleep(time.Millisecond * 500)
			}
		}
	}(ctx)

}

func TestContextDeadline(t *testing.T) {
	// 可以指定一个时间点
	ctx, cancelFunc := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*5))
	// 没有这句话也可以
	defer cancelFunc()
	//_ = cancelFunc

	for {
		select {
		case <-ctx.Done():
			fmt.Println("cancel!")
			return
		default:
			fmt.Println("do word...")
			time.Sleep(time.Second)
		}
	}
}
