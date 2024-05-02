package init2

import (
	"fmt"
	_ "gomod_test/init-test/init1"
)

var Init2 = "look init2"

func init() {
	fmt.Println("init2")

}
