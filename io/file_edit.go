package myio

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// 示例：向 /etc/hosts 中注入域名dns和ip，如果域名已经存在，就修改所有该域名对应的 ip 地址
// 做法：
//		创建 hosts.bak；
//		逐行读取 hosts，并写入 hosts.bak; 判断每行是否包含关键词 dns，如果包含就替换改行内容
// 		最后把 hosts.bak 文件写入 hosts
func InjectionIngress(dns string, ip string) (err error) {
	if ip == "" || dns == "" {
		logs.Info("dns or ip is empty")
		return
	}

	filePath := "/etc/hosts"
	line := ""
	// 首先创建一个空文件 hosts.bak, 如果存在先删除
	err = os.Remove(filePath + ".bak")
	if err == nil || os.IsNotExist(err) {

	} else if err != nil {
		logs.Error(err)
		return
	}
	filebak, err := os.OpenFile(filePath+".bak", os.O_RDWR|os.O_CREATE, 0666) // 读写方式打开文件,文件没有就创建
	if err != nil {
		logs.Error(err)
		return
	}
	defer filebak.Close()

	file, err := os.OpenFile(filePath, os.O_RDWR, 0666) // 读写方式打开文件，将文件内容写到 filebak
	if err != nil {
		logs.Error(err)
		return
	}
	defer file.Close()

	//------------- 读取 file 文件内容到 io 中，加工后，写入 filebak 文件 -------------//
	reader := bufio.NewReader(file)
	replacement := false // 是否发生了行内容替换
	for {
		line, err = reader.ReadString('\n') // 读取每一行内容
		if err != nil {                     // 读到末尾
			if err == io.EOF {
				err = nil
				break
			} else {
				logs.Error(err)
				return
			}
		}
		// 根据关键词替换覆盖当前行
		// for i, dn := range domainNames {
		if strings.Contains(line, dns) {
			str := ip + " " + dns + "\n"
			filebak.WriteString(str)
			replacement = true
		} else {
			filebak.WriteString(line)
		}
	}
	if !replacement {
		filebak.WriteString(ip + " " + dns + "\n")
	}

	//------------- 接下来把 hosts.bak 文件的内容写到原 hosts -------------//
	// Seek设置下一次读/写的位置; 它返回新的偏移量（相对开头）和可能的错误; 参考：https://blog.csdn.net/DisMisPres/article/details/103165731
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		logs.Error(err)
		return
	}
	if _, err = filebak.Seek(0, io.SeekStart); err != nil {
		logs.Error(err)
		return
	}
	err = file.Truncate(0) // 将文件清空
	if err != nil {
		logs.Error(err)
		return
	}
	readerbak := bufio.NewReader(filebak)
	var pos int64 // 偏移量
	for {
		line, err = readerbak.ReadString('\n') //读取每一行内容
		if err != nil {                        //读到末尾
			if err == io.EOF {
				err = nil
				break
			} else {
				logs.Error(err)
				return
			}
		}
		if _, err = file.WriteAt([]byte(line), pos); err != nil {
			logs.Error(err)
			return
		}
		pos += int64(len(line))
	}
	return
}
