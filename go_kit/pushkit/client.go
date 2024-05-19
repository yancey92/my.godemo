package pushkit

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"gitlab.gumpcome.com/common/go_kit/idkit"
	"gitlab.gumpcome.com/common/go_kit/pushkit/gpushkit"
	"gitlab.gumpcome.com/common/go_kit/pushkit/jpushkit"
	"sort"
	"strings"
	"time"
)

var (
	gClient *gpushkit.Client

	jAppKey   string
	jSecret   string
	md5Secret string
)

//InitClient 初始化客户端
func InitClient(conf *PushConfig) error {
	//设置极光推送
	jAppKey = conf.JPushKey
	jSecret = conf.JPushSecret
	md5Secret = conf.Md5Secret

	//设置个推推送内容
	gtAppKey := conf.GtPushKey
	gtSecret := conf.GtPushSecret
	gtAppID := conf.GtAppID

	var err error
	gClient, err = gpushkit.InitClient(gtAppKey, gtSecret, gtAppID)
	if err != nil {
		beego.Error(fmt.Sprintf("个推认证失败， error:%v", err))
	}
	return err
}

//PushMsg 发送push信息, 返回那个通道出错
func PushMsg(msg *PushMessage) error {
	if gClient == nil {
		return fmt.Errorf("未初始化个推服务")
	}

	if msg == nil {
		return fmt.Errorf("push msg消息信息为空")
	}

	errStr := ""
	nErr := 0
	data := createMsg(msg)
	for _, v := range msg.Recv {
		if v.ChanType == JPUSHTYPE {
			data["msg"] = msg.Msg
			err := jPush(v.RegId, data)
			if err != nil {
				beego.Error(fmt.Sprintf("极光推送错误， error:%v, data:%#v", err, data))
				errStr += err.Error()
				nErr++
			}
		} else if v.ChanType == GTPUSHTYPE {
			data["msg"] = msg.Msg
			err := gtPush(v.RegId, data)
			if err != nil {
				beego.Error(fmt.Sprintf("个推推送错误， error:%v, data:%#v", err, data))
				errStr += err.Error()
				nErr++
			}
		}
	}

	if nErr >= len(msg.Recv) {
		return fmt.Errorf(errStr)
	}
	return nil
}

func jPush(regId string, content map[string]interface{}) error {
	buf, err := json.Marshal(content)
	if err != nil {
		beego.Error(fmt.Sprintf("格式化极光心跳的时候出错，buf：%v", buf))
		return err
	}
	cl := jpushkit.NewClient(jAppKey, jSecret)

	msg := jpushkit.NewMessage()
	msg.SetContent(string(buf))

	audience := jpushkit.NewAudience()
	ids := []string{regId}
	audience.RegistrationID = ids

	req := jpushkit.NewPushRequest()
	req.SetPushMsg(msg)
	req.SetAudience(audience)

	m, err := cl.Push(req)
	if err != nil {
		beego.Error(fmt.Sprintf("极光推送出货信息出错， 返回值map:%#v, error:%v", m, err))
		return err
	}

	return nil
}

func gtPush(regId string, content map[string]interface{}) error {
	if gClient == nil {
		return fmt.Errorf("未初始化个推发送客户端信息")
	}

	buf, err := json.Marshal(content)
	if err != nil {
		beego.Error(fmt.Sprintf("格式化极光心跳的时候出错，buf：%v", buf))
		return err
	}
	msg := gpushkit.NewMessage()
	trans := gpushkit.NewTransmission()
	req := gpushkit.NewPushRequest()

	trans.SetContent(string(buf))
	req.SetTransmission(trans)
	req.SetPushMsg(msg)
	req.SetCID(regId)

	m, err := gClient.SafePushSingle(req)
	if err != nil {
		beego.Error(fmt.Sprintf("个推推送出货信息出错， 返回值map:%#v, error:%v", m, err))
		return err
	}

	return nil
}

func createMsg(msg *PushMessage) map[string]interface{} {
	data := make(map[string]interface{})
	data["msg_id"] = idkit.CreateUniqueId()
	data["msg_type"] = msg.MsgType
	data["timeout"] = msg.Timeout
	data["create_time"] = time.Now().UnixNano() / (1000 * 1000)
	data["svm_id"] = msg.SvmId
	signName := signString(data, md5Secret)
	data["sign"] = signName
	return data
}

func signString(data map[string]interface{}, secert string) string {
	buf := bytes.Buffer{}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(fmt.Sprintf("%v", data[k]))
	}

	signStr := buf.String() + secert
	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		beego.Error("在进行MD5加密的时候，取货通知加密失败", err)
	}

	signByte := c.Sum(nil)
	sign := strings.ToUpper(hex.EncodeToString(signByte))
	return sign
}
