package demo

import (
	"fmt"
	"testing"
)

func TestMatchDemo(t *testing.T) {
	boo, err := MatchDemo()
	fmt.Println(boo)
	fmt.Println(err)
}
