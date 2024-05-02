package _interface

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

type Super struct {
	Data interface{} `json:"data"`
}

func test1() {
	var (
		super=&Super{}
		stu =&Student{}
		)
	//stu.Name="zhangsan"
	//stu.Age=23

	err:=json.Unmarshal([]byte(`{"data":{"name":"zhangsan","age":23}}`),super)
	bts,_:=json.Marshal(super.Data)
	err=json.Unmarshal(bts,stu)
	fmt.Println(err,stu)
}