package json_demo_test

import (
	"encoding/json"
	"testing"

	"github.com/yangxinxin/testdemo/json_demo"
)

// 将结构体序列化
func TestStrToJSON(t *testing.T) {
	json_demo.StructSerialization()
}

// 将结构体序列化
// 测试 json 以及 omitempty 的使用
func TestToString(t *testing.T) {
	l := &json_demo.Log{
		Level:    1,
		Filename: "ss111",
	}
	bts, _ := json.Marshal(l)
	t.Logf("%v", string(bts))
}
