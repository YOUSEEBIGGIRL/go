package mytest

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 一个不安全的 sync.Once，不安全的原因是只简单的使用一个 flag + atomic 实现，
// 而标准库的 sync.Once 使用了双重检测
// 下面会演示为什么需要双重检测
type UnsafeOnce struct {
	done uint32
}

func (o *UnsafeOnce) Do(f func()) {
	// CAS 成功才能执行下面的 f()，否则直接返回
	if !atomic.CompareAndSwapUint32(&o.done, 0, 1) {
		return
	}
	f()
}

// 一个需要初始化的对象，未初始化完成，不能使用该对象
type needInitObj struct {
	status int32
}

// 该方法需要 n 初始化完成才能执行
func (n *needInitObj) Run() error {
	if n.status == 1 {
		fmt.Println("obj.run ok")
		return nil
	}
	return errors.New("this obj not init")
}

func TestUnsafe(t *testing.T) {
	var (
		once UnsafeOnce
		wg sync.WaitGroup
		obj needInitObj
	)

	// 一个初始化函数，需要较长时间
	f := func() {
		// 模拟较长的初始化时间
		time.Sleep(time.Second * 5)
		obj.status = 1 	// 初始化完成
		fmt.Println("init done!")
	}

	num := 5
	wg.Add(num)

	for i := 0; i < num; i++ {
		i := i

		go func() {
			defer wg.Done()
			// 第一个 CAS 成功的 goroutine 会执行 f()，但是 f() 需要较长时间，
			// 而此时之后的所有 goroutine 都会因为 CAS 失败而直接 return，继续
			// 执行下面的 obj.Run() ，但是因为此时的 obj 还未初始化完成，所以会
			// 产生一些不可预料的错误
			once.Do(f)
			if err := obj.Run(); err != nil {
				fmt.Printf("goroutine[%d] use not init obj! \n", i)
				return
			}
		}()
	}

	// Output:
	// goroutine[2] use not init obj!
	// goroutine[3] use not init obj!
	// goroutine[1] use not init obj!
	// goroutine[0] use not init obj!
	// init done!
	// obj.run ok

	wg.Wait()
}

// 标准库的 sync.Once 使用了双重检测，不会出现上面的问题
func TestSyncOnce(t *testing.T) {
	var (
		once sync.Once
		wg sync.WaitGroup
		obj needInitObj
	)

	// 一个初始化函数，需要较长时间
	f := func() {
		// 模拟较长的初始化时间
		time.Sleep(time.Second * 5)
		obj.status = 1 	// 初始化完成
		fmt.Println("init done!")
	}

	num := 5
	wg.Add(num)

	for i := 0; i < num; i++ {
		i := i

		go func() {
			defer wg.Done()
			once.Do(f)
			if err := obj.Run(); err != nil {
				fmt.Printf("goroutine[%d] use not init obj! \n", i)
				return
			}
		}()
	}

	// Output:
	// init done!
	// obj.run ok
	// obj.run ok
	// obj.run ok
	// obj.run ok
	// obj.run ok

	wg.Wait()
}
