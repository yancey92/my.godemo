package nsqkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nsqio/go-nsq"
	"gitlab.gumpcome.com/common/go_kit/timekit"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestInitNsqClient(t *testing.T) {
	nsqClient := &NsqClient{}
	nsqConfig := &NsqConfig{}
	nsqClient.Init("123.206.79.254:4761", nsqConfig)
}

func TestSendMessage(t *testing.T) {
	nsqClient := &NsqClient{}
	nsqConfig := &NsqConfig{}
	nsqClient.Init("123.206.79.254:4761", nsqConfig)
	nsqClient.CreateProducer()
	for i := 0; i < 10; i++ {
		nsqClient.SendMessage("write_test", "test"+strconv.Itoa(i))
	}
}

func TestCreateConsumer(t *testing.T) {
	count := 0
	wg := &sync.WaitGroup{}
	wg.Add(1)
	nsqClient := &NsqClient{}
	nsqConfig := &NsqConfig{}
	nsqClient.Init("123.206.79.254:4761", nsqConfig)
	nsqClient.CreateConsumer("write_test", "test", func(message *nsq.Message) error {
		count++
		fmt.Printf("%v-%v\n", count, string(message.Body))
		if string(message.Body) == "test1" {
			return errors.New("测试消息重试机制")
		}
		return nil
	})
	wg.Wait()
}

func HandleJsonMessage(message *nsq.Message) error {
	fmt.Println(string(message.Body))
	return nil
}

type OrderBaseInfo struct {
	SvmId       int    `json:"svm_id" desc:"设备编号"`
	CompanyId   int    `json:"company_id" desc:"公司号"`
	GoodsId     string `json:"goods_id" desc:"商品GID"`
	GoodsName   string `json:"goods_name" desc:"商品名称"`
	GoodsType   int    `json:"goods_type" desc:"商品类型"`
	GoodsTotal  int    `json:"goods_total" desc:"购买商品件数"`
	PayWay      int    `json:"pay_way" desc:"支付方式"`
	OutTradeNo  string `json:"out_trade_no" desc:"商户订单号"`
	TradeNo     string `json:"trade_no" desc:"平台订单号"`
	TradeStatus int    `json:"trade_status" desc:"交易状态"`
	OrderFee    int    `json:"order_fee" desc:"订单金额"`
	PlanType    int    `json:"plan_type" desc:"营销方案类型"`
	PayFee      int    `json:"pay_fee" desc:"支付金额"`
	CreateTime  string `json:"create_time" desc:"订单收录时间"`
	PayUserId   string `json:"pay_user_id" desc:"支付用户ID"`
}

func GetTradeNo(svmId int, payWay int, goodsTotal int) (string, string) {
	_, nowDate, _ := timekit.GetNowTimeMsAndDate(timekit.DateFormat_YYYYMMDDHHMMSSMS)
	_, nowTime, _ := timekit.GetNowTimeMsAndDate(timekit.DateFormat_YYYY_MM_DD_HH_MM_SS)
	return fmt.Sprintf("%v%v%v%v", nowDate, svmId, payWay, goodsTotal), nowTime
}

func TestSendOrderMessage(t *testing.T) {
	nsqClient := &NsqClient{}
	nsqConfig := &NsqConfig{}
	nsqClient.Init("123.206.79.254:4761", nsqConfig)
	nsqClient.CreateProducer()
	svmIds := make([]int, 0, 10000)
	//创建10000个设备
	for i := 20000; i < 25000; i++ {
		svmIds = append(svmIds, i)
	}
	wg := &sync.WaitGroup{}
	totalMoney := make([]int, 0, 10000)
	for _, svmId := range svmIds {
		wg.Add(1)
		go func(id int) {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			payWay := r.Intn(12) + 1
			goodsTotal := r.Intn(1) + 1
			payFee := (r.Intn(10) + 1) * 100
			outTradeNo, createTime := GetTradeNo(id, payWay, goodsTotal)
			orderInfo := OrderBaseInfo{}
			orderInfo.OutTradeNo = outTradeNo
			orderInfo.PayFee = payFee
			totalMoney = append(totalMoney, payFee)
			orderInfo.PayWay = payWay
			orderInfo.CreateTime = createTime
			orderInfo.SvmId = id
			orderInfo.CompanyId = 66666
			fmt.Printf("%v\n", &orderInfo)
			msg, _ := json.Marshal(&orderInfo)
			nsqClient.SendMessage("pay_success", string(msg))
			wg.Done()
		}(svmId)
	}
	wg.Wait()
	temp := 0
	for _, money := range totalMoney {
		temp += money
	}
	fmt.Println(temp)
}
