package flag_demo

import (
	"flag"
	"fmt"
)

var (
	HostIP = ""
	Port   int
)


func Flag() {
	fmt.Println("into function Flag()")
	if flag.Parsed() {
		fmt.Println("flag.Parse() 已经被 go test 命令执行了")
	} else {
		fmt.Println("flag.Parse() 没有被 go test 命令执行")
	}

	flag.StringVar(&HostIP, "host_ip", "127.0.0.1", "宿主机ip，默认空")
	Port = *flag.Int("port", 8080, "服务端口")
}
