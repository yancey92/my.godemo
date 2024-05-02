package init_demo_test

import (
	"flag"
	"fmt"
	"testing"
)

func init() {
	fmt.Println("testing.Init() 被手动显示执行")

	/* Init registers testing flags. These flags are automatically registered by the "go test" command before running test functions,
	so Init is only needed when calling functions such as Benchmark without using "go test". */
	testing.Init()

	if flag.Parsed() {
		fmt.Printf("%v\n", "flag.Parse() 被 testing.Init() 执行了")
	} else {
		fmt.Printf("%v\n", "flag.Parse() 没有被 testing.Init() 执行")
	}
}

/*
	该测试验证了目前的 go 版本中 testing.Init() 中不会执行 flag.Parse()
	但是在进入测试函数之前，flag.Parse() 会被 go test 命令执行。

	另外：testing.Init() 会在进入测试函数之前被 go test 命令执行；不使用"go test"时，如在调用 Benchmark 等函数，testing.Init()不被自动执行
*/
func TestInit(t *testing.T) {
	t.Logf("%v", "测试函数开始")
	if flag.Parsed() {
		t.Logf("%v", "在进入测试函数前执行了 flag.Parse()")
	} else {
		t.Logf("%v", "在进入测试函数前没有执行 flag.Parse()")
	}
	t.Logf("%v", "测试 testing.Init() 完成")
}
