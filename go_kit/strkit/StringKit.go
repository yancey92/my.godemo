// @Title 字符串处理工具

package strkit

import (
	"bytes"
	"strconv"
	"unicode/utf8"
)

// @Title 判断多个字符串是否不为空
// @Description
// @param strs
// usage:
//	StrNotBlank("a", "15") return true
func StrNotBlank(strs ...string) bool {
	if len(strs) == 0 {
		return false
	}
	for _, v := range strs {
		if v == "" {
			return false
		}
	}
	return true
}

// @Title 判断多个字符串是否为空
// @Description
// @param strs
// usage:
//	StrIsBlank("a", "") return false
func StrIsBlank(strs ...string) bool {
	if len(strs) == 0 {
		return false
	}
	for _, v := range strs {
		if v != "" {
			return false
		}
	}
	return true
}

// @Title Int64 string to int64
func StrToInt64(str string) (int64, error) {
	v, err := strconv.ParseInt(str, 10, 64)
	return int64(v), err
}

// @Title Int string to int
func StrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// @Title 多个字符串拼接
// @Description
// @param strs
// usage:
//	StrJoin("hello ", "world", "", " ", "is go write")
// 	return hello world is go write
func StrJoin(strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	var strBuffer bytes.Buffer
	for _, v := range strs {
		if StrNotBlank(v) {
			strBuffer.WriteString(v)
		}
	}
	return strBuffer.String()
}

// @Title 获取字符串长度
// @Description
// @param str
// usage:
//	GetStrLen("hello ")
// 	return 6
func GetStrLen(str string) int {
	return utf8.RuneCountInString(str)
}

// @Title 字符串构建对象
// @Description
// usage:
//	实例构造方法: StringBuilder{}
type StringBuilder struct {
	buf bytes.Buffer
}

// @Title 添加字符串到字符串构建实例里面
// @Description 空字符串将会被忽略
// usage:
//	实例构造方法: sb.Append("hello").Append(" world")
func (sb *StringBuilder) Append(str string) *StringBuilder {
	if StrNotBlank(str) {
		sb.buf.WriteString(str)
	}
	return sb
}

// @Title 输出字符串构建实例里面的所有字符串
// @Description 空字符串将会被忽略
// usage:
//	实例构造方法: sb.Append("hello").Append(" world").ToString()
func (sb *StringBuilder) ToString() string {
	return sb.buf.String()
}

// @Title 获取字符串长度
// @Description
// @param str ,start开始下标，结束下标（包含）
//	SubStr("20170620120101", 0,6)
// 	return 201706
func SubStr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if end < 0 || end > length {
		return ""
	}
	return string(rs[start:end])
}
