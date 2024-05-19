// @Title 业务错误码
package logiccode

import (
	"github.com/astaxie/beego"
	"gitlab.gumpcome.com/common/go_kit/strkit"
	"strconv"
)

func New(code int, msg string) error {
	return &LogicCode{code, msg}
}

// @Title 业务错误码
// @Description
//	业务错误码由6位组成,前3位代表类别,后3位代表具体业务。
//	100XXX: 100代表通用类别;
// 		XXX<=100: 	代表DAO层错误;
// 		100<XXX<=200: 	代表logic层错误;
// 		200<XXX<=300: 	代表controller层错误;
type LogicCode struct {
	Code int    `json:"code" desc:"业务错误码"`
	Msg  string `json:"msg" desc:"错误描述"`
}

func (code *LogicCode) Error() string {
	return strkit.StrJoin("[", strconv.Itoa(code.Code), "]<", code.Msg, ">")
}

// 获取error的状态码
func GetCode(err error) int {
	if err == nil {
		return 0
	}
	switch value := err.(type) {
	case *LogicCode:
		return value.Code
	default:
		beego.Info(value)
		return 0
	}
	return 0
}

// @Title DB连接错误
// @Description 用于DAO层操作DB错误反馈
func DbConErrorCode() error {
	return New(100001, "db connect error")
}

// @Title DB插入操作错误
// @Description 用于DAO层操作DB错误反馈
func DbInsertErrorCode() error {
	return New(100002, "db insert error")
}

// @Title DB更新操作错误
// @Description 用于DAO层操作DB错误反馈
func DbUpdateErrorCode() error {
	return New(100003, "db update error")
}

// @Title 通过主键ID更新DB操作错误
// @Description 用于DAO层操作DB错误反馈
func DbUpdateByIdErrorCode() error {
	return New(100004, "db update by id is nil")
}

// @Title DB删除操作错误
// @Description 用于DAO层操作DB错误反馈
func DbDeleteErrorCode() error {
	return New(100005, "db delete error")
}

// @Title DB查询操作错误
// @Description 用于DAO层操作DB错误反馈
func DbQueryErrorCode() error {
	return New(100006, "db query error")
}

// @Title DB分页查询超出总页数范围
// @Description 用于DAO层操作DB错误反馈
//	页码错误、每页显示记录总数错误。
func DbPageOutErrorCode() error {
	return New(100007, "db page query out of range")
}

// @Title DB操作影响记录数为0
// @Description 用于DAO层插入、更新、删除记录时没有实际发生影响记录数
func DbZeroErrorCode() error {
	return New(100008, "db affected rows is 0")
}

// @Title 记录总数值字符串转整形异常
// @Description
func DbPageCountToIntCode() error {
	return New(100009, "page count string to int error")
}

// @Title 数据库配置名称错误
// @Description
func DbConfigNameErrorCode() error {
	return New(100010, "mysql config name is nil")
}

// @Title DB结果字符串转整型错误
// @Description 用于DAO层操作DB错误反馈
func DbItemToIntErrorCode() error {
	return New(100011, "db item to int error")
}

// @Title 连接Mongo数据库错误
func MongoConnErrorCode() error {
	return New(100012, "mongo conn error")
}

// @Title Mongo Session 克隆错误
func MongoSessionCloneErrorCode() error {
	return New(100013, "mongo seesion clone error")
}

// @Title Mongo Session 初始化错误
func MongoSessionErrorCode() error {
	return New(100014, "mongo is not inited or session is nil")
}

// @Title Mongo 搜索条件为空
func MongoParamsErrorCode() error {
	return New(100015, "mongo params is nil")
}

// @Title Mongo 添加或者更新操作错误
func MongoUpsertErrorCode(err error) error {
	return New(100016, err.Error())
}

// @Title Mongo 删除操作错误
func MongoRemoveErrorCode(err error) error {
	return New(100017, err.Error())
}

// @Title Redis Key 不存在
func RedisKeyErrorCode() error {
	return New(100018, "redis get key does not exists")
}

// @Title Redis client 错误
func RedisClientErrorCode() error {
	return New(100019, "redis client is nil")
}

// @Title Redis 参数错误
func RedisParamsErrorCode() error {
	return New(100020, "redis params is empty")
}

// @Title 请求参数值错误
// @Description 用于反射请求参数对象、参数值类型转换、必填参数校验错误反馈
func ReqParamErrorCode() error {
	return New(100301, "param value error")
}
