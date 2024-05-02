package goto_demo

import "fmt"

// 会死循环
func Demo1() {
Err:
	fmt.Println("End")
Success:
	fmt.Println("End")

	fmt.Println("begin")
	var err error
	fmt.Println(err)

	if err != nil {
		goto Err
	} else {
		goto Success
	}
}

// 
func Demo2() {
	fmt.Println("begin")
	var err error
	fmt.Println(err)

	if err != nil {
		goto Err
	} else {
		goto Success
	}

Err:
	fmt.Println("End Err")
Success:
	fmt.Println("End Success")
}
