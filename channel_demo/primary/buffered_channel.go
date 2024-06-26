package primary

//------------------------------------------------------------------------------------------------------------------------

// 往channel中发送数据，如果channel是没有缓存的，那么发一条数据到channel，该channel就会阻塞，
// 等待goroutine来取出里面的数据，这样就会导致channel性能降低。
// 所以，带缓冲的channel就比较好的解决了channel每往里面发一条数据就会导致channel阻塞的问题。
// 下面这段代码，是不会报错的，。
func bufferedChannel1() {
	c := make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	//c <- 4// 如果这一行代码放开，就会报错"deadlock"
}

//------------------------------------------------------------------------------------------------------------------------

