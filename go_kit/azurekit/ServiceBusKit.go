package azurekit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"gitlab.gumpcome.com/common/go_kit/strkit"
	"gitlab.gumpcome.com/common/go_kit/timekit"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ServiceBusLogDevMode  = "dev"
	ServiceBusLogProdMode = "prod"
)

// 服务总线消息模型
type ServiceBusMsgMode struct {
	XMLName xml.Name `xml:"entry"`
	Content string   `xml:"content"`
}

type myServiceBus struct {
	logger *logs.BeeLogger
	isInit bool
}

var mySerBus = new(myServiceBus)

// @Title 初始化服务总线日志
// @param mode 日志模式 dev 开发模式 prod 生产模式
func initServiceBusLogger(mode string) {
	logDir := "servicebus_logs"
	logFile := logDir + string(os.PathSeparator) + "servicebus.log"
	os.MkdirAll(logDir, 0777)

	mySerBus.logger = logs.NewLogger()
	mySerBus.logger.SetLogger(logs.AdapterMultiFile, `{
		"filename":"`+logFile+`",
		"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]
	}`)
	mySerBus.logger.Async(1e3) //1000

	if ServiceBusLogProdMode == mode {
		mySerBus.logger.SetLevel(beego.LevelInformational)
	} else {
		mySerBus.logger.SetLevel(beego.LevelDebug)
	}
}

// @Title 获取服务总线令牌
func getServiceBusToken(uri string, keyName string, key string) string {
	if strkit.StrIsBlank(uri, keyName, key) {
		mySerBus.logger.Error("uri, keyName, key is blank!")
		return ""
	}

	nowSs, _, _ := timekit.GetNowTimeSsAndDate(timekit.DateFormat_YYYY_MM_DD)
	//7天后token失效
	expiry := strconv.FormatInt(nowSs+60*60*24*7, 10)

	stringToSign := url.QueryEscape(uri) + "\n" + expiry
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return "SharedAccessSignature sr=" + url.QueryEscape(uri) + "&sig=" + url.QueryEscape(signature) + "&se=" + expiry + "&skn=" + keyName
}

// @Title 初始化服务总线
// @param logMode 日志模式 dev 开发模式 prod 生产模式
func InitServiceBus(logMode string) {
	if mySerBus.isInit {
		mySerBus.logger.Info("服务总线已初始化完成,无需重复初始化...")
	}
	initServiceBusLogger(logMode)
	mySerBus.isInit = true
	mySerBus.logger.Info("服务总线初始化成功...")
}

// @Title 发送消息到服务总线中
// @param uri 		服务总线URI
// @param keyName 	密钥名称
// @param key 		密钥
// @param message 	消息
func SendMessageToServiceBus(uri string, keyName string, key string, message string) error {
	token := getServiceBusToken(uri, keyName, key)
	if strkit.StrIsBlank(token, message) {
		return errors.New("token or message is blank!")
	}

	req := httplib.Post(uri + "/messages")
	req.Header("Authorization", token)
	req.Header("ContentType", "application/atom+xml;type=entry;charset=utf-8")
	req.SetTimeout(60*time.Second, 60*time.Second)

	msgBuffer := strkit.StringBuilder{}
	msgBuffer.Append("<entry xmlns='http://www.w3.org/2005/Atom'><content type='application/xml'>").Append(message).Append("</content></entry>")
	req.Body(msgBuffer.ToString())
	_, err := req.String()
	resp, err := req.Response()
	if err != nil {
		mySerBus.logger.Error("发送消息到服务总线请求失败, uri=%s, msg=%s, err=%v", uri, message, err)
		return errors.New("send message to queue request is error!")
	}
	if http.StatusCreated != resp.StatusCode {
		return errors.New(fmt.Sprintf("send message to queue response status code is %d", resp.StatusCode))
	}
	return nil
}

// @Title 从服务总线接收消息
func receiveMessageFromServiceBus(uri string, keyName string, key string, subName string, msgProcessor func(msg string) bool) {
	token := getServiceBusToken(uri, keyName, key)
	if strkit.StrIsBlank(token) {
		mySerBus.logger.Error("token is blank!")
	}

	mySerBus.logger.Info(fmt.Sprintf("开始接收服务总线消息, uri=%s", uri))
	req := httplib.Post(uri + "/subscriptions" + subName + "/messages/head")
	req.Header("Authorization", token)
	req.SetTimeout(60*time.Second, 60*time.Second)

	result, err := req.String()
	resp, err := req.Response()
	if err != nil {
		mySerBus.logger.Error("从队列接收消息请求失败, uri=%s, err=%v", uri, err)
	}
	if http.StatusCreated != resp.StatusCode || result == "" {
		return
	}
	brokerProperties := resp.Header.Get("BrokerProperties")
	location := resp.Header.Get("Location")
	fmt.Println(brokerProperties)

	//解析消息
	msg := ServiceBusMsgMode{}
	err = xml.Unmarshal([]byte(result), &msg)
	if err != nil {
		mySerBus.logger.Error("从服务总线接收消息后解析消息失败, uri=%s, msg=%s, err=%v", uri, result, err)
	}
	if msg.Content == "" {
		mySerBus.logger.Error("从服务总线接收消息后解析消息内容为空, uri=%s, msg=%s, err=%v", uri, result, err)
	}
	isSuccess := msgProcessor(strings.TrimSpace(msg.Content))
	if isSuccess {
		mySerBus.logger.Info(fmt.Sprintf("开始请求删除服务总线消息 location=%s", location))
		req = httplib.Delete(location)
		req.Header("Authorization", token)
		req.SetTimeout(60*time.Second, 60*time.Second)
		_, err = req.String()
		resp, err = req.Response()
		if err != nil {
			mySerBus.logger.Error("删除服务总线中消息请求失败, location=%s, err=%v", location, err)
		}
		if http.StatusOK != resp.StatusCode {
			mySerBus.logger.Error("删除服务总线中消息处理失败, location=%s, err=%v", location, err)
		}
	}
}

// @Title 启用接收服务总线消息服务
// @param uri 			服务总线URI
// @param keyName 		密钥名称
// @param key 			密钥
// @param subName		订阅名称
// @param func(msg string) bool 消息处理器
func StartReceiveMessageFromServiceBusServer(uri string, keyName string, key string, subName string, msgProcessor func(msg string) bool) {
	go func() {
		for {
			receiveMessageFromServiceBus(uri, keyName, key, subName, msgProcessor)
		}
	}()
}
