package primary

import "fmt"

// channel 是可以被close的（通常在发送方来close），这样就通知接收方"我没有要发送的数据了"，
// 接收方也能因此而判断停止接收数据。
// 当channel被close后，接收方依旧可以从channel中取数据，此时取出的数据是发送数据类型的"零值"。


//------------------------------------------------------------------------------------------------------------------------
// 下面这个代码，从channel中取数据时没有做任何判断，当channel被close后，如果goroutine依旧在运行，
// 那么goroutine会依旧从channel中取数据，此时取出的数据是发送数据类型的"零值"
func closeChannel1() {
	c:=make(chan int)
	go worker_close1(0,c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	close(c)
	//time.Sleep(time.Millisecond)
}


func worker_close1(id int ,c chan int) {
	for  {
		fmt.Printf("worker id is %d, worker received %c\n",id,<-c)
	}
}

//------------------------------------------------------------------------------------------------------------------------
// 下面这个代码，当channel被close后，goroutine就不从channel中取数据了。
func closeChannel2() {
	c:=make(chan int)
	go worker_close2(0,c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	close(c)
	//time.Sleep(time.Millisecond)
}

func worker_close2(id int ,c chan int) {
	for  {
		if n,ok:=<-c;ok{
			fmt.Printf("worker id is %d, worker received %c\n",id,n)
		}else {
			break
		}
	}
}

//------------------------------------------------------------------------------------------------------------------------
// 从channel中取数据是可以通过 range 来遍历的，这样就不用像上面那样从 channel 中取数据时还要判断是否取到数据
func rangeChannel() {
	c:=make(chan int)
	go worker_range(0,c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	close(c)
	//time.Sleep(time.Millisecond)
}

func worker_range(id int ,c chan int) {
	for n := range c {
		fmt.Printf("worker id is %d, worker received %c\n",id,n)
	}
}
//------------------------------------------------------------------------------------------------------------------------

