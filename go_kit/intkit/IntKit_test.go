package intkit

import (
	"fmt"
	"testing"
)

func TestIntIsZero(t *testing.T) {
	if !IntIsZero(0, 0) { //应该是true
		t.Log("TestIntIsZero-01 fail")
		t.Fail()
	}

	if IntIsZero(0, 1) { //应该是true
		t.Log("TestIntIsZero-01 fail")
		t.Fail()
	}
}

func TestToFixed(t *testing.T) {
	f := ToFixed(0.12345, 6)
	fmt.Println(f)
}
