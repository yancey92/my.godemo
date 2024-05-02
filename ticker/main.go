package main

import (
	"fmt"
	"time"
)

func main() {

	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {
		fmt.Println("aaaa")
		ticker.Stop()
	}

	fmt.Println("bbbbb")
}
