package dbkit

import (
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego"
	"gitlab.gumpcome.com/common/go_kit/logiccode"
	"gitlab.gumpcome.com/common/go_kit/strkit"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
)

const (
	SOCKET_TIME_OUT_MS = 10000 //操作超时,默认10s
	CONN_TIME_OUT_MS   = 10000 //连接超时,默认10s
	MAX_POOL_SIZE      = 100   //连接池大小
	MAX_RETRIES        = 5     //连接失败后,重试次数
)

var (
	mongoInited   bool //是否已初始化
	globalSession *mgo.Session
	connUrl       string //完整连接URL
	mongoDBName   string
)

type MongoSearch struct {
	Collection string
	Key        string
	Value      interface{}
}

// 非SSL协议初始K化数据库
// @connUrl 连接字符串
// @dbName  数据库名称
// @maxConn 最大连接数
func InitMongoDB(connUrl string, dbName string, maxConn int) {
	if connUrl == "" || dbName == "" {
		panic("conn url or db name is empty!")
	}
	if mongoInited {
		return
	}
	fullUrl := setConnUrlOptions(connUrl)
	mySession, err := mgo.Dial(fullUrl)
	if err != nil {
		panic(err)
	}
	mySession.SetMode(mgo.Monotonic, true)
	if maxConn == 0 {
		mySession.SetPoolLimit(MAX_POOL_SIZE)
	} else {
		mySession.SetPoolLimit(maxConn)
	}
	globalSession = mySession
	connUrl = fullUrl
	mongoDBName = dbName
	mongoInited = true
}

// SSL协议初始化数据库
// @connUrl 连接字符串
// @dbName  数据库名称
func InitMongoDBWithSSL(connUrl string, dbName string) {
	if connUrl == "" || dbName == "" {
		panic("conn url or db name is empty!")
	}
	if mongoInited {
		return
	}
	fullUrl := setConnUrlOptions(connUrl)
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	dialInfo, err := mgo.ParseURL(fullUrl)
	if err != nil {
		panic(err)
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	mySession, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	err = mySession.Ping()
	if err != nil {
		panic(err)
	}
	mySession.SetMode(mgo.Monotonic, true)
	mySession.SetPoolLimit(100)
	globalSession = mySession
	connUrl = fullUrl
	mongoDBName = dbName
	mongoInited = true
}

// 查询出与search匹配的结果,如果没有就添加,否则会覆盖原有记录
// 注意:如果Mongo搜索到多条与search匹配的记录,只会更新最新插入的一条记录。
func MongoUpsertDoc(search *MongoSearch, doc interface{}) (*mgo.ChangeInfo, bool, error) {
	session, err := getSession()
	if err != nil {
		return &mgo.ChangeInfo{}, false, err
	}
	defer session.Close()
	if search == nil || doc == nil || search.Collection == "" || search.Key == "" || search.Value == nil {
		return &mgo.ChangeInfo{}, false, logiccode.MongoParamsErrorCode()
	}
	changeInfo, err := session.DB(mongoDBName).C(search.Collection).Upsert(bson.M{search.Key: search.Value}, doc)
	if err != nil {
		return &mgo.ChangeInfo{}, false, logiccode.MongoUpsertErrorCode(err)
	}
	return changeInfo, true, nil
}

//插入记录
func MongoInsert(colelection string, data interface{}) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB(mongoDBName).C(colelection)
	err = c.Insert(data)
	if err != nil {
		return err
	}
	return nil
}

//查找记录处理器
func MongoFindHandler(collection string, fun func(*mgo.Collection)) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB(mongoDBName).C(collection)
	fun(c)
	return nil
}

// 删除所有记录
func MongoRemoveAllDoc(search *MongoSearch) (bool, error) {
	session, err := getSession()
	if err != nil {
		return false, err
	}
	defer session.Close()
	if search == nil || search.Collection == "" || search.Key == "" || search.Value == "" {
		return false, logiccode.MongoParamsErrorCode()
	}
	_, err = session.DB(mongoDBName).C(search.Collection).RemoveAll(bson.M{search.Key: search.Value})
	if err != nil {
		return false, logiccode.MongoRemoveErrorCode(err)
	}
	return true, nil
}

//查找记录
func MongoFindDoc(search *MongoSearch) (*mgo.Query, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	data := session.DB(mongoDBName).C(search.Collection).Find(bson.M{search.Key: search.Value})
	return data, nil
}

// 设置连接字符串后缀可选项
func setConnUrlOptions(connUlr string) string {
	opts := make([]string, 0)
	opts = append(opts, connUlr)
	opts = append(opts, "?")
	opts = append(opts, "authMechanism=MONGODB-CR")
	opts = append(opts, "&maxPoolSize=100")
	//opts = append(opts, "&connectTimeoutMS=10000") //10s连接超时
	//opts = append(opts, "&socketTimeoutMS=10000")  //10s操作超时
	return strkit.StrJoin(opts...)
}

// 获取Mongo连接
func GetMongoSession() (*mgo.Session, error) {
	return getSession()
}

// 获取集合
func GetCollection(colelection string) (*mgo.Collection, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}

	//defer session.Close()
	c := session.DB(mongoDBName).C(colelection)
	return c, nil
}

func getSession() (*mgo.Session, error) {
	if !mongoInited || globalSession == nil {
		return nil, logiccode.MongoSessionErrorCode()
	}
	isSessionOk := true
	globalSession.Refresh()
	err := globalSession.Ping()
	if err != nil {
		beego.Error(fmt.Sprintf("globalSession ping fail %v", err))
		isSessionOk = false
		globalSession.Refresh()
		for i := 0; i < MAX_RETRIES; i++ {
			err = globalSession.Ping()
			if err == nil {
				isSessionOk = true
				beego.Info("Reconnect to mongodb successful.")
				break
			} else {
				beego.Error(fmt.Sprintf("Reconnect to mongodb fail:%v"), i)
			}
		}
	}
	if isSessionOk {
		return globalSession.Clone(), nil
	}
	return nil, logiccode.MongoSessionCloneErrorCode()
}
