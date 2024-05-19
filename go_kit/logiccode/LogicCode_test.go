package logiccode

import (
	"fmt"
	"testing"
)

func TestLogicCode_Error(t *testing.T) {
	err := New(100001, "error message 错误描述")
	fmt.Printf("%v", err)
}
