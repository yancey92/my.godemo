package pushkit

//PushConfig 配置信息
type PushConfig struct {
	JPushKey     string
	JPushSecret  string
	GtPushKey    string
	GtPushSecret string
	GtAppID      string
	Md5Secret    string
}

type Receiver struct {
	RegId    string
	ChanType int
}

//PushMessage 要发送的推送消息
type PushMessage struct {
	Timeout int
	MsgType int
	SvmId   int
	Msg     interface{}
	Recv    []Receiver
}

var (
	JPUSHTYPE  = 10
	GTPUSHTYPE = 11
)
