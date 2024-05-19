package azurekit

import (
	"fmt"
	"testing"
	"time"
)

func TestSendMessageToServiceBus(t *testing.T) {
	InitServiceBus("dev")

	uri := "https://dyltest.servicebus.chinacloudapi.cn/dyltest-topic"
	keyName := "RootManageSharedAccessKey"
	key := "iEZvPvwcSFe+U4gOFyZx7F4s2tIrznyNvMzn/4tpuAI="
	message := "你好, Azure ServiceBus"

	err := SendMessageToServiceBus(uri, keyName, key, message)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println("消息发送到队列成功...")
}

func TestStartReceiveMessageFromServiceBusServer(t *testing.T) {
	InitServiceBus("dev")

	uri := "https://dyltest.servicebus.chinacloudapi.cn/dyltest-topic"
	keyName := "RootManageSharedAccessKey"
	key := "iEZvPvwcSFe+U4gOFyZx7F4s2tIrznyNvMzn/4tpuAI="
	subName := "sub01"

	msgProcessor := func(msg string) bool {
		fmt.Println("开始处理消息")
		time.Sleep(2 * time.Second)
		fmt.Println("msgProcessor " + msg)
		fmt.Println("消息处理结束")
		return true
	}

	StartReceiveMessageFromServiceBusServer(uri, keyName, key, subName, msgProcessor)
	time.Sleep(5 * time.Minute)
}
