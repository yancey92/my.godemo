package string_test

import (
	"fmt"
	"testing"
)

func TestStringPoint(t *testing.T) {
	name := new(string)
	t.Logf("name:%v\n", name)
	t.Logf("*name:%v\n", *name)
	t.Logf("&name:%v\n", &name)

	fmt.Println("-------------------------------------------")
	*name = "zhangsan"
	t.Logf("name:%v\n", name)
	t.Logf("*name:%v\n", *name)
	t.Logf("&name:%v\n", &name)
}
