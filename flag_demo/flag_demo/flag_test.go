package flag_demo_test

import (
	"flag"
	"testing"

	"github.com/yangxinxin/testdemo/flag_demo/flag_demo"
	// _ "github.com/yangxinxin/testdemo/flag_demo/flag_demo"
)

// 经过测试，发现：flag.Parse()  会被 go test 命令执行
func TestFlag(t *testing.T) {
	flag_demo.Flag()
	if !flag.Parsed() {
		flag.Parse()
	}

	t.Logf("HostIP is:%v", flag_demo.HostIP)
	t.Logf("Port is:%v", flag_demo.Port)
}
