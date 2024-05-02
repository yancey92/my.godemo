package myslice

import (
	"fmt"
	"testing"
)

func TestDemo1(t *testing.T) {
	x := make([]int, 0, 10)
	x = append(x, 1, 2, 3)
	y := append(x, 4)
	z := append(x, 5)

	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(z)

	// %p	address of 0th element in base 16 notation, with leading 0x
	// https://pkg.go.dev/fmt
	fmt.Printf("%p\n", x)
	fmt.Printf("%p\n", y)
	fmt.Printf("%p\n", z)
}

func TestDemo2(t *testing.T) {
	a := make([]int, 0, 1)
	fmt.Printf("%p\n", a)
	a = append(a, 1)

	fmt.Printf("%p，%p\n", a, &a[0])

	a = append(a, 2)
	fmt.Printf("%p，%p\n", a, &a[0])
}
