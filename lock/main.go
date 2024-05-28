package main

import (
	"fmt"
	"sync"
)

var str = ""
var rwMutex sync.RWMutex

// 如果 str是空，就设置一个值；如果 str不为空，就直接返回
func getStr() string {
	rwMutex.RLock()
	if str != "" {
		defer rwMutex.RUnlock()
		return str
	}
	rwMutex.RUnlock()

	rwMutex.Lock()
	defer rwMutex.Unlock()
	if str == "" {
		setStr("hello world")
	}

	return str
}

func setStr(s string) {
	//rwMutex.Lock() // 这个地方不能再加锁了，加了锁就会导致死锁
	//defer rwMutex.Unlock()
	str = s
}

func main() {
	s := getStr()
	fmt.Printf(s)
}
