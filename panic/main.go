package main

import (
	"fmt"
	"time"
)

func main() {

	go func() {
		fmt.Println("a")
	}()

	go func() {
		fmt.Println("b")
		panic("b")
	}()
	time.Sleep(3 * time.Second)
	fmt.Println("c")
}
