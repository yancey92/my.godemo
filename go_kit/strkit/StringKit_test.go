package strkit

import (
	"testing"
)

func TestStrNotBlank(t *testing.T) {
	if !StrNotBlank("a") { //应该是true
		t.Log("TestStrNotBlank-01 fail")
		t.Fail()
	}

	if !StrNotBlank("a", "15") { //应该是true
		t.Log("TestStrNotBlank-02 fail")
		t.Fail()
	}

	if StrNotBlank("a", "", "     ") { //应该是false
		t.Log("TestStrNotBlank-03 fail")
		t.Fail()
	}

	if StrNotBlank("") { //应该是false
		t.Log("TestStrNotBlank-04 fail")
		t.Fail()
	}

	if StrNotBlank() { //应该是false
		t.Log("TestStrNotBlank-05 fail")
		t.Fail()
	}

	if !StrNotBlank("   ") { //应该是true
		t.Log("TestStrNotBlank-06 fail")
		t.Fail()
	}
}

func TestStrIsBlank(t *testing.T) {
	if StrIsBlank("a") { //应该是false
		t.Log("TestStrIsBlank-01 fail")
		t.Fail()
	}

	if StrIsBlank("a", "15") { //应该是false
		t.Log("TestStrIsBlank-02 fail")
		t.Fail()
	}

	if StrIsBlank("a", "", "     ") { //应该是false
		t.Log("TestStrIsBlank-03 fail")
		t.Fail()
	}

	if StrIsBlank("", "   ") { //应该是false
		t.Log("TestStrIsBlank-04 fail")
		t.Fail()
	}

	if StrIsBlank() { //应该是false
		t.Log("TestStrIsBlank-05 fail")
		t.Fail()
	}

	if !StrIsBlank("", "") { //应该是true
		t.Log("TestStrIsBlank-06 fail")
		t.Fail()
	}
}

func TestStrJoin(t *testing.T) {
	if !("hello world is go write" == StrJoin("hello ", "world", "", " ", "is go write")) {
		t.Log("TestStrJoin-01 fail")
		t.Fail()
	}
}

func TestStringBuilder_ToString(t *testing.T) {
	sb := StringBuilder{}
	content := sb.Append("hello").Append(" world is ").Append("go write").ToString()
	if !(content == StrJoin("hello ", "world", "", " ", "is go write")) {
		t.Log("TestStrJoin-01 fail")
		t.Fail()
	}
}

func TestGetStrLen(t *testing.T) {
	enLen := GetStrLen("hello ")
	if 6 != enLen {
		t.Fail()
	}

	cnLen := GetStrLen("你好,Go")
	if 5 != cnLen {
		t.Fail()
	}
}
