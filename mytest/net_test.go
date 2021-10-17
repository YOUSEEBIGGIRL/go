package mytest

import (
	"net"
	"testing"
)

func TestNetListen(t *testing.T) {
	l, _ := net.Listen("tcp", ":8080")	
	c, _ := l.Accept()	
	c.Read(nil)
}