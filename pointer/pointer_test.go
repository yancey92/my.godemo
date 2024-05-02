package pointer

import (
	"fmt"
	"testing"
)

func TestDemo1(t *testing.T) {
	a := 1
	fmt.Printf("%p\n", &a)
	b := &a
	fmt.Printf("%p\n", b)
	fmt.Printf("%v\n", *b)
	c := &b
	fmt.Printf("%p\n", c)
}
