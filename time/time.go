package my_time

import "time"

/*
	The time zone database needed by LoadLocation may not be present on all systems, especially non-Unix systems.
	LoadLocation looks in the directory or uncompressed zip file named by the ZONEINFO environment variable,
	if any, then looks in known installation locations on Unix systems,
	and finally looks in $GOROOT/lib/time/zoneinfo.zip.

	LoadLocation 所需的时区数据库可能并非存在于所有系统上，尤其是非 Unix 系统。
	LoadLocation 首先查找由 ZONEINFO 环境变量（如果有）命名的目录或未压缩的 zip 文件，
	然后在 Unix 系统上的已知安装位置中查找 (通常是 /usr/share/zoneinfo )，
	最后查看 $GOROOT/lib/time/zoneinfo.zip。
*/
// 获取当前时间对应时区 Asia/Chongqing 的时间
func TimeLocalFormat() string {
	format := "2006-01-02 15:04:05"
	nowtime := time.Now().UTC()
	timelocal, _ := time.LoadLocation("Asia/Chongqing")
	return nowtime.In(timelocal).Format(format)
}


