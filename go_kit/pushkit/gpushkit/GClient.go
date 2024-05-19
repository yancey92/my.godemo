package gpushkit

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	client *Client
)

//Client 个推推送
type Client struct {
	AppKey string
	Secret string
	AppID  string
	GtAuth string
	Mutex  *sync.RWMutex
}

func getTimeMilliInt() int64 {
	return int64(time.Now().UnixNano() / (1000 * 1000))
}
func getTimeMilliString() string {
	return strconv.FormatInt(getTimeMilliInt(), 10)
}

//InitClient 初始化信息, 只能初始化一次非线程安全的
func InitClient(appKey, secret, appID string) (*Client, error) {
	client = &Client{AppKey: appKey, Secret: secret, AppID: appID}
	client.Mutex = new(sync.RWMutex)
	//首先需要做的就是要完成认证,个推的认证
	auth, err := AuthSign(appKey, secret, appID)
	if err != nil {
		beego.Error(fmt.Sprintf("个推认证失败， error:%v", err))
		return client, err
	}

	client.Mutex.Lock()
	client.GtAuth = auth
	client.Mutex.Unlock()
	return client, nil
}

//SafePushSingle 线程安全
func (c *Client) SafePushSingle(req *PushRequest) (map[string]interface{}, error) {

	req.PushMsg.SetAppKey(c.AppKey)

	c.Mutex.RLock()
	m, err := PushSingle(c.AppID, c.GtAuth, req)
	c.Mutex.RUnlock()
	if err != nil {
		beego.Error(fmt.Sprintf("个推push信息失败，error:%v,返回值%#v", err, m))
	}

	result := m["result"]
	strRes, ok := result.(string)
	if !ok || strRes != CONSTANT_OK {
		//再次进行权限验证
		auth, err := AuthSign(c.AppKey, c.Secret, c.AppID)

		if err != nil {
			beego.Error(fmt.Sprintf("个推再次认证失败， error:%v", err))
			return m, err
		}

		c.Mutex.Lock()
		c.GtAuth = auth
		c.Mutex.Unlock()

		c.Mutex.RLock()
		next, err := PushSingle(c.AppID, c.GtAuth, req)
		c.Mutex.RUnlock()

		if err != nil {
			beego.Error(fmt.Sprintf("个推再次push信息出错, %#v, error:%v", next, err))
			return next, err
		}

		res := next["result"]
		strRes, ok := res.(string)
		if !ok || strRes != CONSTANT_OK {
			beego.Error(fmt.Sprintf("个推再次push信息出错,result结果不是ok, %#v, error:%v", next, err))
			return next, err
		}
	}
	return m, err
}

//AuthSign 鉴权信息
func AuthSign(appKey, masterSecret, appID string) (string, error) {
	req := &SignRequest{}

	req.Timestamp = getTimeMilliString()
	h := sha256.New()
	h.Write([]byte(appKey + req.Timestamp + masterSecret))
	md := h.Sum(nil)
	req.Sign = hex.EncodeToString(md)
	req.AppKey = appKey

	url := CONSTANT_PUSH_URL + appID + "/auth_sign"
	buf, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := request("POST", url, bytes.NewReader(buf), "")
	if err != nil {
		return "", err
	}

	mapResponse, err2 := resp.Map()
	if err2 != nil {
		return "", err2
	}

	result := mapResponse["result"]
	strRes, ok := result.(string)
	if !ok || strRes != CONSTANT_OK {
		return "", fmt.Errorf("请求失败，返回值result不是ok,response:%v, 参数appkey:%s, masterSecret:%s, appID:%s", mapResponse, appKey, masterSecret, appID)
	}

	appAuth := mapResponse["auth_token"]
	strRes, ok = appAuth.(string)
	if !ok {
		return "", fmt.Errorf("鉴权失败，未取到auth_token值,response:%v, 参数appkey:%s, masterSecret:%s, appID:%s", mapResponse, appKey, masterSecret, appID)
	}
	return strRes, nil
}

//PushSingle 发送单个push信息
func PushSingle(appID string, auth string, req *PushRequest) (map[string]interface{}, error) {
	url := CONSTANT_PUSH_URL + appID + "/push_single"
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := request("POST", url, bytes.NewReader(buf), auth)
	if err != nil {
		return nil, err
	}
	return resp.Map()
}

//PushList 设置群发送
func PushList(appID string, auth string, cIDs []string, req *SaveMsgsRequest) (map[string]interface{}, error) {
	url := CONSTANT_PUSH_URL + appID + "/save_list_body"
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := request("POST", url, bytes.NewReader(buf), auth)
	if err != nil {
		return nil, err
	}

	mapResponse, err := resp.Map()
	if err != nil {
		return nil, err
	}

	result := mapResponse["result"]
	strRes, ok := result.(string)
	if !ok || strRes != CONSTANT_OK {
		return nil, fmt.Errorf("请求失败,返回值result不是ok,response:%v,参数appID:%s,authtoken:%s", mapResponse, appID, auth)
	}

	taskID, ok := mapResponse["task_id"]
	if !ok {
		return nil, fmt.Errorf("没有找到返回的任务id,response:%v,参数appID:%s,authtoken:%s", mapResponse, appID, auth)
	}

	strTask, ok := taskID.(string)
	if !ok {
		return nil, fmt.Errorf("格式化任务标识出错:%v,response:%v,appID:%s,authtoken:%s", taskID, mapResponse, appID, auth)
	}

	pListReq := &PushListRequest{CID: cIDs, TaskID: strTask}
	urlList := CONSTANT_PUSH_URL + appID + "/push_list"

	bufList, err := json.Marshal(pListReq)
	if err != nil {
		return nil, err
	}
	resp, err = request("POST", urlList, bytes.NewReader(bufList), auth)
	if err != nil {
		return nil, err
	}
	return resp.Map()
}

func request(method, url string, body io.Reader, auth string) (*Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Charset", "UTF-8")
	req.Header.Set("Content-Type", "application/json")

	if len(auth) > 0 {
		req.Header["authtoken"] = []string{auth}
	}

	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{data: buf}, nil
}
