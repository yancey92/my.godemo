package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"my.swagger/model"
	"net/http"
	"time"
)

// 拦截器
// 设置允许跨域的处理器
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置允许所有来源的跨域请求
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		// 可以根据需要设置其他跨域相关的头，例如：
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// 允许浏览器在跨域请求中携带Cookie
		// 同时服务端要设置为：Access-Control-Allow-Origin 就不能设置为*，而必须指定确切的、信任的源域名，以确保安全性。
		// 且客户端也要做适配（哎，放弃了、、、）
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8001")
		// 如果是预检请求（OPTIONS请求），直接返回204状态码
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

//  swag fmt;
// 	swag init -g ./cmd/main.go -o ./swag_static --generatedTime --ot "yaml"
//
//	swagger general API info
//	@title			swagger demo API
//	@version		v1.0
//	@description	swagger demo
//	@termsOfService

//	@contact.name	yancey
//	@contact.url
//	@contact.email	yangxinxin_mail@163.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@schemes	http
//	@host		localhost:8000
//	@BasePath	/
//	@Accept		json
//	@Produce	json

// @externalDocs.description
// @externalDocs.url
func main() {

	mux := http.NewServeMux()

	// 设定静态文件目录
	staticDir := http.Dir("./swag_static")
	// 创建一个用于服务静态文件的 FileServer
	fileServer := http.FileServer(staticDir)
	// 将 "/swagger" 路径下的所有请求都代理给 fileServer 处理
	// 这意味着如果你的静态文件目录中有 "index.html"，那么通过 "/swagger/index.html" 就可以访问到
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fileServer))
	mux.HandleFunc("/api/v1/cookieVerify", cookieVerify)

	// 使用中间件包装 ServeMux
	handler := corsMiddleware(mux)
	server := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      handler,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
		return
	}
}

// ---
//
//	@Summary		访问首页
//	@Description	这是一个测试接口，测试服务是否就绪1
//	@Description	这是一个测试接口，测试服务是否就绪2
//	@Tags			swagger
//	@Param			requestBody		body		model.ReqBody	true	"人员信息"
//	@Param			X-Project-Id	header		string			true	"Project id"
//	@Param			Cookie			header		string			false	"cookie"	<----"swagger ui 中虽然让你填写Cookie，而实际上不生效"
//
//	@Success		200				{object}	string			"请求发送成功"
//	@Failure		500				{object}	string			"程序内部错误"
//	@Router			/api/v1/cookieVerify [post]
//
// swagger ui 本质上是不支持传递 Cookie 的，即便你生成的 swagger 文档中有 Cookie 相关的设置。
// 使用中如果你需要使用 Cookie，可以在浏览器打开 swagger ui，并设置对应的 Cookie，前提是swagger-ui和后端服务必须是同源，
// 这就需要将 swagger-ui 一同打包，作为后端服务的一部分。
func cookieVerify(w http.ResponseWriter, r *http.Request) {
	var (
		requestBody = &model.ReqBody{}
	)

	reqBodyBt, err := io.ReadAll(r.Body)
	logrus.Infof("receive request body:%v", string(reqBodyBt))
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(reqBodyBt, requestBody); err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	cks := r.Cookies()
	projectId := r.Header.Get("X-Project-Id")

	fmt.Printf("X-Project-Id:%v\n", projectId)
	fmt.Printf("cookie:%v\n", cks)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, "hello, "+fmt.Sprintf("%#v", requestBody))
}
