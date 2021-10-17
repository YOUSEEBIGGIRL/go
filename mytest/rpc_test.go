package mytest

import (
	"log"
	"net/http"
	"net/rpc"
	"testing"
)

type Cal struct {
	X, Y int64
}

type Res struct {
	R int64
}

func (c *Cal) Add(param *Cal, res *Res) error {
	res.R = param.X + param.Y
	return nil
}

func TestServer(t *testing.T) {
	if err := rpc.Register(new(Cal)); err != nil {
		log.Fatalln(err)
	}
	rpc.HandleHTTP()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}

func TestClient(t *testing.T) {
	dial, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}

	cal := new(Cal)
	cal.X = 100
	cal.Y = 50

	res := new(Res)

	//go func() {
	//	call := dial.Go("Cal.Add", cal, res, nil)
	//	select {
	//	case c := <-call.Done:
	//		log.Println(c.Reply)
	//	case <-time.After(time.Second * 5):
	//		log.Println("timeout")
	//	}
	//}()
	//time.Sleep(time.Second * 10)


	if err := dial.Call("Cal.Add", cal, res); err != nil {
		log.Println(err)
		return
	}

	log.Println(res)
}
