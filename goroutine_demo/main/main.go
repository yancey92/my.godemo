package main

import (
	"errors"
	"fmt"
	"time"
)

func a(times int) (err error) {
	fmt.Printf("a, times :%v\n", times)
	return
}

func b(times int) (err error) {
	fmt.Printf("b, times :%v\n", times)
	return
}
func c(times int) (err error) {
	fmt.Printf("c, times :%v\n", times)
	err = errors.New("err")
	return
}

func loopFunction(end chan bool, f func(i int) (err error)) {
	go func() {
		times := 0
		for {
			times++
			select {
			case <-end:
				return
			default:

			}
			if err := f(times); err != nil {
				time.Sleep(1 * time.Second)
			} else {
				return
			}

		}

	}()

}

// 循环调用某个方法，如果该方法返回err，则停止调用
func main() {
	end := make(chan bool, 0)
	for i := 0; i < 3; i++ {
		loopFunction(end, c)
		loopFunction(end, a)
	}

	time.Sleep(5 * time.Second)
	close(end)

}
