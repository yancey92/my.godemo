package jpushkit

import "encoding/json"

//Audience 发送的人
type Audience struct {
	RegistrationID []string `json:"registration_id,omitempty"`
}

//Message 设置message信息
type Message struct {
	Content     string                 `json:"msg_content"`
	Title       string                 `json:"title,omitempty"`
	ContentType string                 `json:"content_type,omitempty"`
	Extras      map[string]interface{} `json:"extras,omitempty"`
}

//PushRequest 这个是请求
type PushRequest struct {
	ID       string    `json:"cid,omitempty"`
	Platform string    `json:"platform"`
	PushMsg  *Message  `json:"message,omitempty"`
	Audience *Audience `json:"audience,omitempty"`
}

//NewPushRequest 新建一个推送请求
func NewPushRequest() *PushRequest {
	req := PushRequest{
		Platform: "android",
	}
	return &req
}

//NewMessage 新建一个透传消息
func NewMessage() *Message {
	m := Message{
		Title:       "default",
		Content:     "甘来推送",
		ContentType: "text",
	}
	return &m
}

//NewAudience 新建一个受众
func NewAudience() *Audience {
	au := Audience{}
	return &au
}

//SetPlatform 设置平台信息
func (req *PushRequest) SetPlatform(platform string) {
	req.Platform = platform
}

//SetPushMsg 设置push信息
func (req *PushRequest) SetPushMsg(msg *Message) {
	req.PushMsg = msg
}

//SetAudience 设置接收人
func (req *PushRequest) SetAudience(audience *Audience) {
	req.Audience = audience
}

//SetContent 消息内容
func (m *Message) SetContent(content string) {
	m.Content = content
}

//SetTitle 标题
func (m *Message) SetTitle(title string) {
	m.Title = title
}

//SetContentType ...
func (m *Message) SetContentType(strType string) {
	m.ContentType = strType
}

//AddExtras 加入扩展字段
func (m *Message) AddExtras(key string, value interface{}) {
	if m.Extras == nil {
		m.Extras = make(map[string]interface{})
	}
	m.Extras[key] = value
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
