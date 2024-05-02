package json_demo

import (
	"encoding/json"
	"fmt"
)

type Log struct {
	Level    int    `json:"level,omitempty" `
	Filename string `json:"filename,omitempty"`
	Color    bool   `json:"color,omitempty"`
	Daily    bool   `json:"daily,omitempty"`
	Maxdays  int    `json:"maxdays,omitempty"`
	Rotate   bool   `json:"rotate"`
}

// 结构体 serialization
func StructSerialization() {
	type Demo struct {
		Command string `json:"command"`
	}

	d := Demo{}
	d.Command = `curl --location --request PUT 'http://dashboard.lenovo.com/api/v1/registering/register' \
	--header 'Content-Type: application/json' \
	--data-raw '{
		"datacenter_name": "rancher-test",
		"token": "MGMCAQACEQDS52MK0Nw1QIfRfSaTOhdnAgMBAAECEQCV/r6v6I9Uxv/J3tc5onvBAgkA6o1Gly9yqVECCQDmMIZkMe6HNwIIaKE68yhMWvECCBXOAuz6zd2BAgkAlR9YVQCWm2A=",
		"cmo_addr": "10.121.112.49:30009",
		"comment":"rancher 测试，不要随意删除：10.121.118.123",
		"longitude":"11",
		"latitude":"22",
		"edge_type":"KubeEdge"
	}'`

	bt, _ := json.Marshal(d) // 将结构体序列化
	fmt.Println(string(bt))
}
