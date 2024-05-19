// @Title 唯一ID生成工具
package idkit

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"
)

// @Title 创建32位唯一ID
// @Description
// usage:
//	CreateUniqueId()
func CreateUniqueId() string {
	return strings.ToUpper(NewV1().String())
}

// @Title 创建字符串MD5
// @Description
// usage:
//	CreateMd5("abc")
func CreateMd5(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}
