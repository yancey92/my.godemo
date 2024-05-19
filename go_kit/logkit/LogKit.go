package logkit

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"os"
)

const (
	LogDevMode  = "dev"
	LogProdMode = "prod"
)

// 输出错误日志信息
type ErrorInfo struct {
	RequestId     string                 `desc:"请求唯一标识"`
	MethodName    string                 `desc:"方法名"`
	MethodContext *context.Context       `desc:"请求参数信息"`
	ExtContext    map[string]interface{} `desc:"拓展请求参数信息"`
	ErrorRemark   string                 `desc:"错误信息描述"`
	ErrorMsg      error                  `desc:"错误信息"`
}

func InitLog() {
	logmode := beego.AppConfig.String("logmode")
	if logmode == "" || (logmode != LogDevMode && logmode != LogProdMode) {
		panic("config logmode is empty or log mode is not dev or prod!")
	}

	logDir := "logs"
	logFile := logDir + string(os.PathSeparator) + "server.log"
	os.MkdirAll(logDir, 0777)

	if logmode == LogDevMode {
		beego.SetLevel(beego.LevelDebug)
	} else {
		beego.SetLevel(beego.LevelInformational)
	}

	beego.SetLogger(logs.AdapterMultiFile, `{
		"filename":"`+logFile+`",
		"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]
	}`)
}

func (this *ErrorInfo) AddExtContent(key string, val interface{}) *ErrorInfo {
	if this.ExtContext == nil {
		this.ExtContext = make(map[string]interface{})
	}
	this.ExtContext[key] = val
	return this
}

func OutErrorInfo(errorMsg *ErrorInfo) string {
	bodyContent := ""
	formContent := ""
	methodContent := ""
	if errorMsg.MethodContext != nil {
		bodyContent = string(errorMsg.MethodContext.Input.RequestBody)
		formContent = fmt.Sprintf("%#v", errorMsg.MethodContext.Request.Form)
		methodContent = errorMsg.MethodContext.Input.Context.Request.RequestURI
	}
	if methodContent == "" {
		methodContent = errorMsg.MethodName
	}
	return fmt.Sprintf("RequestId:%v MethodName:%v RequestBody:%v RequestForm:%v ExtContext:%#v ErrorRemark:%v ErrorMsg:%v",
		errorMsg.RequestId, methodContent, bodyContent,
		formContent, errorMsg.ExtContext, errorMsg.ErrorRemark, errorMsg.ErrorMsg)
}
