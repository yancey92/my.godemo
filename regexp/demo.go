package demo

import "regexp"

func MatchDemo() (bool, error) {
	// regexp.MatchString(pattern string, s string)   pattern 为正则表达式，s 为需要校验的字符
	// 匹配ipv4地址
	return regexp.MatchString(`dash`, "dashborad.len")
}
