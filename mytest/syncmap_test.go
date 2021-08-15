package mytest

import (
	"sync"
	"testing"
)

func TestName(t *testing.T) {
	var cmap sync.Map

	cmap.Store("k1", "v1")
	cmap.Store("k2", "v2")
	cmap.Load("k1")
	cmap.Load("k2")
	cmap.Store("k3", "v3")
	cmap.Delete("k3")
}
