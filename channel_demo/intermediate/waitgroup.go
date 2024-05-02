package intermediate

import (
	"fmt"
	"sync"
)

//------------------------------------------------------------------------------------------------------------------------
// sync.WaitGroup 对象内部有一个计数器，最初从0开始，它有三个方法：Add(), Done(), Wait() 用来控制计数器的数量。
// Add(n) 把计数器的值加n ；Done() 把计数器-1 ；wait() 会阻塞代码的运行，直到计数器地值减为0才会放行。

type worker struct {
	in chan int
	done func() // 抽象出来一个方法，具体实现在 create worker 时配置
}

func channel1() {
	var wg sync.WaitGroup
	woker:=createWorker1(0,&wg)
	wg.Add(20) // 计数器在原来0的基础上加20，因为下面要发送20个字符

	go func() {
		for i:=0;i<10 ;i++  {
			woker.in <- i
		}
	}()
	go func() {
		for i:=10;i<20 ;i++  {
			woker.in <- i
		}
	}()

	wg.Wait() // 阻塞代码，防止程序退出，直到计数器为0才放行
}

func createWorker1(id  int,wg *sync.WaitGroup) worker {
	w:=worker{
		in: make(chan int),
		done: func() {
			wg.Done() // 具体要做的动作
		},
	}
	go worker1(id,w)
	return w
}

func worker1(id int ,w worker) {
	for n:= range w.in {
		fmt.Printf("worker id is %d, worker received %d\n",id,n)
		w.done() // 每完成一次任务，计数器减一
	}
}
//------------------------------------------------------------------------------------------------------------------------
