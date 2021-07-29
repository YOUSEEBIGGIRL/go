package mytest

import (
	"fmt"
	"testing"
)

func TestMake(t *testing.T) {
	m := make(map[int]int, 123)
	m[1] = 2
	fmt.Println()
}

func Test123(t *testing.T) {
	fmt.Println("123")
}


