package primary

import "fmt"

//------------------------------------------------------------------------------------------------------------------------

// 下面报错：fatal error: all goroutines are asleep - deadlock!
// "c <- 1" 这一行报的错，原因是 channel 是 goroutine 与 goroutine 之间的交互通道，
// 程序发现这里只往 channel 中发送数据，却没有从中取数据，所以报错
func channel1() {
	c := make(chan int)
	c <- 1
	c <- 2
	n := <-c
	fmt.Println(n)
}
//------------------------------------------------------------------------------------------------------------------------

// 下面这个代码有可能输出1，也有可能输出1和2
// 原因是：该channel是无缓存的的，当向channel中发送数据1，channel就会等着goroutine从它里面取数据（阻塞），
// 		将1取出来，然后进2，取2时程序有可能还没有来得及打印程序就退出了
func channel2() {
	c := make(chan int)
	go func() {
		for {
			n := <-c
			fmt.Println(n)
		}
	}()
	c <- 1
	c <- 2
}
//------------------------------------------------------------------------------------------------------------------------

func channel3() {
	// 创建channel
	var chans [10]chan int
	for i:=0;i<10 ; i++ {
		chans[i]=createWorker(i)
	}
	// 给channel分发数据
	for i:=0;i<10 ;i++  {
		chans[i]<-'a'+i
	}
	//time.Sleep(time.Millisecond)
}

func createWorker(id  int) chan int {
	c:=make(chan int)  //创建完channel，当然是要让goroutine使用了，所以下面紧接着就把goroutine使用的函数写出来（这个思想很好）
	go worker(id,c)
	return c
}

func worker(id int ,c chan int) {
	for  {
		fmt.Printf("worker id is %d, worker received %c\n",id,<-c)
	}
}
//------------------------------------------------------------------------------------------------------------------------


