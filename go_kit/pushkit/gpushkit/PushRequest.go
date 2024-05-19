package gpushkit

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
	//"gitlab.gumpcome.com/common/go_kit/idkit"
)

//Message 设置消息类型
type Message struct {
	AppKey  string `json:"appkey"`
	Offline bool   `json:"is_offline"`
	MsgType string `json:"msgtype"`
}

//Style 模式
type Style struct {
	StyleType int    `json:"type"`
	Text      string `json:"text"`
	Title     string `json:"title"`
}

//Transmission 透传信息
type Transmission struct {
	TransType    bool   `json:"transmission_type"`
	TransContent string `json:"transmission_content"`
}

//PushRequest 推送请求
type PushRequest struct {
	CID          string        `json:"cid,omitempty"`
	RequestID    string        `json:"requestid,omitempty"`
	Transmission *Transmission `json:"transmission"`
	PushMsg      *Message      `json:"message,omitempty"`
}

//SaveMsgsRequest 保存消息信息
type SaveMsgsRequest struct {
	PushMsg      *Message      `json:"message,omitempty"`
	Transmission *Transmission `json:"transmission"`
}

//PushListRequest 发送群发消息信息
type PushListRequest struct {
	CID    []string `json:"cid,omitempty"`
	TaskID string   `json:"taskid"`
}

//SignRequest 鉴权信息
type SignRequest struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	AppKey    string `json:"appkey"`
}

//NewPushRequest 建立推送信息请求
func NewPushRequest() *PushRequest {
	f := func() string {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(10000)
		return getTimeMilliString() + strconv.Itoa(n)
	}
	req := PushRequest{RequestID: f()}

	return &req
}

//NewSaveMsgsRequest 建立新的一个请求
func NewSaveMsgsRequest() *SaveMsgsRequest {
	req := &SaveMsgsRequest{}
	return req
}

//NewTransmission 建立透传消息
func NewTransmission() *Transmission {
	trans := &Transmission{
		TransType:    false,
		TransContent: "",
	}
	return trans
}

//NewMessage 建立消息类型
func NewMessage() *Message {
	m := &Message{Offline: false, MsgType: "transmission"}
	return m
}

//SetPushMsg 设置消息类型
func (req *PushRequest) SetPushMsg(m *Message) {
	req.PushMsg = m
}

//SetTransmission 设置透传
func (req *PushRequest) SetTransmission(trans *Transmission) {
	req.Transmission = trans
}

//SetCID 设置cid
func (req *PushRequest) SetCID(cid string) {
	req.CID = cid
}

//SetPushMsg 设置群发送
func (req *SaveMsgsRequest) SetPushMsg(m *Message) {
	req.PushMsg = m
}

//SetTransmission 设置透传信息
func (req *SaveMsgsRequest) SetTransmission(m *Transmission) {
	req.Transmission = m
}

//SetAppKey 设置appkey
func (m *Message) SetAppKey(appKey string) {
	m.AppKey = appKey
}

//SetContent 设置透传内容
func (trans *Transmission) SetContent(content string) {
	trans.TransContent = content
}

//Response 返回数据
type Response struct {
	data []byte
}

//Map 结果数据
func (res *Response) Map() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal(res.data, &result)
	return result, err
}

//Bytes 获取data数据
func (res *Response) Bytes() []byte {
	return res.data
}
