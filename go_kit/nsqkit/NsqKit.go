package nsqkit

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/nsqio/go-nsq"
	"math/rand"
	"time"
)

// 返回可用消息服务结果集
type nsqLookupNodesModel struct {
	Producers []*nsqProducer `json:"producers"`
}

type nsqProducer struct {
	RemoteAddress    string `json:"remote_address"`
	HostName         string `json:"hostname"`
	BroadcastAddress string `json:"broadcast_address"`
	TcpPort          int    `json:"tcp_port"`
	HttpPort         int    `json:"http_port"`
	Version          string `json:"version"`
}

type NsqConfig struct {
	MsgTimeout   time.Duration //消息推送到消费者超时时间,默认30s
	DialTimeout  time.Duration //连接nsqd超时时间,默认10s
	WriteTimeout time.Duration //往nsqd中写消息超时时间,默认5s
	MaxInFlight  int           //消费者单次消费消息最大数量,默认40
}

type NsqClient struct {
	nsqConfig     *NsqConfig
	nsqNodes      *nsqLookupNodesModel
	nsqProducers  []*nsq.Producer
	nsqConsumer   *nsq.Consumer
	nsqLookupAddr string
}

// 初始化消息服务客户端
// @nsqLookupAddr nsqLookup 客户端的地址
func (this *NsqClient) Init(nsqLookupAddr string, nsqConfig *NsqConfig) {
	if nsqLookupAddr == "" {
		beego.Error("NsqLookup address is empty!")
		panic("NsqLookup address is empty!")
	}
	if nsqConfig == nil {
		beego.Error("NsqClient config is nil!")
		panic("NsqClient config is nil!")
	}

	reqHttp := httplib.Get("http://" + nsqLookupAddr + "/nodes")
	model := nsqLookupNodesModel{}
	if err := reqHttp.ToJSON(&model); err != nil {
		errorInfo := fmt.Sprintf("NsqClient get nodes is error! %v\n", err)
		beego.Error(errorInfo)
		panic(errorInfo)
	}
	this.nsqNodes = &model
	this.nsqConfig = nsqConfig
	this.nsqLookupAddr = nsqLookupAddr
	if this.nsqConfig.DialTimeout == 0 {
		this.nsqConfig.DialTimeout = 10 * time.Second //10s
	}
	if this.nsqConfig.MsgTimeout == 0 {
		this.nsqConfig.MsgTimeout = 30 * time.Second //30s
	}
	if this.nsqConfig.WriteTimeout == 0 {
		this.nsqConfig.WriteTimeout = 5 * time.Second //5s
	}
	if this.nsqConfig.MaxInFlight == 0 {
		this.nsqConfig.MaxInFlight = 40
	}
}

// 创建消息生产者
func (this *NsqClient) CreateProducer() {
	config := nsq.NewConfig()
	config.DialTimeout = this.nsqConfig.DialTimeout
	config.MsgTimeout = this.nsqConfig.MsgTimeout
	config.WriteTimeout = this.nsqConfig.WriteTimeout
	this.nsqProducers = make([]*nsq.Producer, 0)
	for _, node := range this.nsqNodes.Producers {
		addr := fmt.Sprintf("%v:%v", node.BroadcastAddress, node.TcpPort)
		prder, _ := nsq.NewProducer(addr, config)
		this.nsqProducers = append(this.nsqProducers, prder)
	}
}

// 发送消息到消息服务
// @topicName 消息服务主题名称
// @message 往主题发送的内容
func (this *NsqClient) SendMessage(topicName string, message string) bool {
	//随机其中任意一个nsqd节点生产消息
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(this.nsqNodes.Producers))
	nsqPrd := this.nsqProducers[idx]
	err := nsqPrd.Ping()
	if err != nil {
		beego.Error(fmt.Sprintf("ping【IP:%v,Topic:%v】消息服务失败,error=%v,msq=%v\n", this.nsqProducers[idx], topicName, err, message))
		return false
	}
	err = this.nsqProducers[idx].Publish(topicName, []byte(message))
	if err != nil {
		beego.Error(fmt.Sprintf("发送消息到【IP:%v,Topic:%v】消息服务失败,error=%v,msq=%v\n", this.nsqProducers[idx], topicName, err, message))
		return false
	}
	return true
}

// 创建消息消费者
// @topicName 消费的消息主题
// @channelName 消费的消息通道
// @handle 消费者处理器
func (this *NsqClient) CreateConsumer(topicName string, channelName string, handle func(message *nsq.Message) error) {
	config := nsq.NewConfig()
	config.DialTimeout = this.nsqConfig.DialTimeout
	config.MsgTimeout = this.nsqConfig.MsgTimeout
	//消费者同时能够处理多少个nsqd的消息
	config.MaxInFlight = this.nsqConfig.MaxInFlight
	consumer, _ := nsq.NewConsumer(topicName, channelName, config)
	consumer.AddHandler(nsq.HandlerFunc(handle))
	err := consumer.ConnectToNSQLookupd(this.nsqLookupAddr)
	if err != nil {
		beego.Error(fmt.Sprintf("connect lookupd fail %v\n", err))
	}
}
