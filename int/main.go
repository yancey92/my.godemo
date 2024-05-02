package main

import "fmt"

func main() {
	fmt.Println(int(^uint32(0)>>1))

}
