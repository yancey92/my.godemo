package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	//test1()
	//test2()
	//test3()
	test4()
}

// 执行下面这个函数有可能什么都不会输出，函数就推出了。
func test1() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				fmt.Printf("hello from goroutine %d\n", i)
			}
		}(i)
	}
}

// 下面这个函数会启动10个goroutine，每个 goroutine 不停的输出，且运行1毫秒退出
// 由于协程中的逻辑是打印一条内容，即IO操作，IO操作会有等待的过程，所以会调用
func test2() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				fmt.Printf("hello from goroutine %d\n", i)
			}
		}(i)
	}
	time.Sleep(time.Millisecond) // sleep 1 ms
}

// 下面这个程序和上面的很像，只是goroutine中不在调用IO，
// 在 go1.13以及之前的版本中，下面的代码中 a[i]++ 会修改内存中的变量，goroutine 是不会主动交出控制权的,
//		即协程之间没有机会切换，所以程序会一直停在第一个协程中不停的死循环执行a[i]++，最后导致计算机cpu飙升，程序卡死
//
// 在 go1.14中 goroutine 开始支持抢占式切换，那么下面的这个代码会有多个协程抢占go调度器，最终运行1毫秒输出结果：
//		Running in go1.14
//		[2870269 3533919 3855788 2150833 0 0 0 0 0 3946632]
//
// 其实，下面这个代码还是有错误的，使用 go run -race 就可以看到数据引用冲突的错误：
// 原因就是，最后一行一边在打印组a，同时又有多个go协程往数组a里面写，这个问题要通过channel来解决
func test3() {
	var a [10]int
	fmt.Println("Running in", runtime.Version())
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				a[i]++
				runtime.Gosched() // 在go1.13以及之前的版本中加上这句话，让 goroutine 主动交出控制权，这样每个协程都有执行的机会
			}
		}(i)
	}
	time.Sleep(time.Millisecond) // sleep 1 ms
	fmt.Println(a)
}

// 运行该函数，我们在控制台 top，看到该进程"活动的线程数"最大为8，因为我的电脑是4核8线程的。
// 这里开了1000个goroutine，但是实际上会映射到我电脑cpu的8个线程上来执行。
func test4() {
	for i := 0; i < 1000; i++ {
		go func(i int) {
			for {
				fmt.Printf("hello from goroutine %d\n", i)
			}
		}(i)
	}
	time.Sleep(time.Minute) // sleep 1 s
}