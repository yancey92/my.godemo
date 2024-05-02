package primary

import "fmt"

// 默认情况下，channel是双向的：两端都可以发数据和取数据
// channel 是可以定义方向的：
// 		var ch1 chan<- float64 对于定义者而言，该channel是只能写入数据的
// 		var ch2 <-chan float64 对于定义者而言，该channel是只能读出数据的
//不能将单向 channel 转换为双向 channel

//------------------------------------------------------------------------------------------------------------------------
func vectorChannel1() {
	ch:=make(chan int)
	go worker_vector1(0,ch)
	ch <- 'a'
	ch <- 'b'
	ch <- 'c'
	close(ch)
}

func worker_vector1(id int,c <-chan int) {
	for n := range c {
		fmt.Printf("worker id is %d, worker received %c\n",id,n)
	}
}
//------------------------------------------------------------------------------------------------------------------------
