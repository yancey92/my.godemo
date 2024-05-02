package main2

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/register", Handler)
	http.ListenAndServe("localhost:8080", nil) // 监听在本机所有网卡的 8080 端口
}

func Handler(rw http.ResponseWriter, r *http.Request) {
	var (
		watch string
	)
	params := strings.Split(r.URL.RawQuery, "&")
	for _, v := range params {
		if strings.HasPrefix(v, "watch=") {
			tmpList := strings.Split(v, "=")
			watch = tmpList[1]
			break
		}
	}
	if watch == "1" {
		rw.Header().Set("Transfer-Encoding", "chunked")
		rw.Header().Set("Access-Control-Allow-Origin","*")
		rw.WriteHeader(200)
	}
	//rw.WriteHeader(200)
	for i := 1; i <= 5; i++ {

		rw.Write([]byte(fmt.Sprintf("%x",4) + "\r\n很多时候我们写的asp程序会因为做很多操作，所以会花上一分钟甚至几分钟时间。为了使软件使用者能够耐心的等待程序的执行，我们经常会希望有一个进度条来表示程序执行的状态。或者最起码要显示一个类似： “数据载入中”，“正在保存数据” 等的说明性文字。此时我们就会用到Response.flush()。他会将缓冲区中编译完成的数据先发送到客户端\r\n")) // 这里是每次要发给客户端的数据
	
		rw.(http.Flusher).Flush()  // 这一步也很重要
		time.Sleep(time.Second * 4)
	}
	rw.Write([]byte("0\r\n\r\n")) // 数据发送结束

}
