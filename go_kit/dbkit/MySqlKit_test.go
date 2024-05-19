package dbkit

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"gitlab.gumpcome.com/common/go_kit/idkit"
	"gitlab.gumpcome.com/common/go_kit/timekit"
	"testing"
)

//INSERT INTO user(name,age,email,gender,height,interests) VALUES (?,?,?,?,?,?)
func TestCreateMysqlInsertSQL(t *testing.T) {
	return
	tableName := "user"
	data := make(map[string]interface{})
	data["name"] = "小刘"
	data["age"] = 18
	data["email"] = "xiaoliu@gumpcome.com"
	data["gender"] = 1
	data["height"] = 180
	data["interests"] = "游泳,爬山"

	sql, params := CreateMysqlInsertSQL(tableName, data)
	fmt.Println(sql)
	fmt.Println(params)
}

var dbConfigName = "main"

func TestInitMysql(t *testing.T) {
	logger := logs.NewLogger(1000)
	logger.SetLogger(logs.AdapterConsole)
	dbUserName := "root"
	dbUserPwd := "dyl123"
	dbHost := "127.0.0.1:3306"
	dbName := "godb"
	dbMaxIdle := 10
	dbMaxActive := 20

	InitMysql(dbUserName, dbUserPwd, dbHost, dbName, dbConfigName, dbMaxIdle, dbMaxActive)
}

func GetConn() *sql.DB {
	conn, err := GetMysqlCon(dbConfigName)
	if err != nil {
		return nil
	}
	return conn
}

func TestSaveInMysql(t *testing.T) {
	return
	tableName := "user"
	data := make(map[string]interface{})
	data["name"] = "大龙"
	data["age"] = 18
	data["email"] = "xiaoliu@gumpcome.com"
	data["gender"] = 1
	data["height"] = 180
	data["interests"] = "游泳,爬山"
	unixTime, createTime, _ := timekit.GetNowTimeMsAndDate(timekit.DateFormat_YYYY_MM_DD_HH_MM_SS)
	data["createtime"] = createTime
	data["unixtime"] = unixTime

	result, id, err := SaveInMysql(GetConn(), tableName, data)
	if !result && err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	fmt.Printf("保存记录返回的主键结果 id=%d\n", id)
}

func TestUpdateByIdInMysql(t *testing.T) {
	return
	tableName := "user"
	data := make(map[string]interface{})
	data["id"] = 12
	data["name"] = "小刘3"
	data["age"] = 18
	data["email"] = "xiaoliu@gumpcome.com"
	data["gender"] = 1
	data["height"] = 180
	data["interests"] = "游泳1,爬山1"
	unixTime, createTime, _ := timekit.GetNowTimeMsAndDate(timekit.DateFormat_YYYY_MM_DD_HH_MM_SS)
	data["createtime"] = createTime
	data["unixtime"] = unixTime

	result, err := UpdateByIdInMysql(GetConn(), tableName, data)
	if !result && err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	fmt.Printf("根据ID更新记录返回结果 result=%t\n", result)
}

func TestUpdateInMysql(t *testing.T) {
	return
	sql := `UPDATE user SET name = ? , age = ? , email = ? , gender = ? , interests = ? , createtime = ? , unixtime = ? , height = ?  WHERE id = ?`
	unixTime, createTime, _ := timekit.GetNowTimeMsAndDate(timekit.DateFormat_YYYY_MM_DD_HH_MM_SS)
	result, err := UpdateInMysql(GetConn(), sql, "大龙", 20, "dalong@gumpcome.com", 1, "潜水,旅游", createTime, unixTime, 176, 13)
	if !result && err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	fmt.Printf("更新记录返回结果 result=%t\n", result)
}

func TestDeleteByIdInMysql(t *testing.T) {
	return
	result, err := DeleteByIdInMysql(GetConn(), "user", 14)
	if !result && err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	fmt.Printf("根据ID删除记录返回结果 result=%t\n", result)
}

func TestDeleteInMysql(t *testing.T) {
	return
	sql := `DELETE FROM user WHERE name = ?`
	result, err := DeleteInMysql(GetConn(), sql, "小刘2")
	if !result && err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	fmt.Printf("记录返回结果 result=%t\n", result)
}

func TestFindInMysql(t *testing.T) {
	//sql := `SELECT name AS user_name, age As user_age FROM user LIMIT 1`
	sql := `SELECT t1.id, t1.age, t2.content, t2.unixtime FROM user t1 LEFT JOIN user_address t2 ON t2.userid = t1.id;`
	intItems := []string{"user_age", "unixtime"}
	result, err := FindInMysql(GetConn(), sql, intItems)
	if err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
}

func TestFindFirstInMysql(t *testing.T) {
	sql := `SELECT name AS user_name, age As user_age FROM user WHERE name LIKE ? AND id < ? LIMIT 1, 10`
	intItems := []string{"user_age"}
	result, err := FindFirstInMysql(GetConn(), sql, intItems, "%黄月英%", 100)
	if err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
}

func TestFindFirstInMysql2(t *testing.T) {
	sql := `SELECT COUNT(*) AS count FROM user WHERE name LIKE ? ORDER BY id`
	result, err := FindFirstInMysql(GetConn(), sql, []string{"count"}, "%大龙%")
	if err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
}

func TestPaginateInMysql(t *testing.T) {
	selectSql := `SELECT name AS user_name, age As user_age`
	sqlExceptSelect := `FROM user WHERE name LIKE ?`
	result, err := PaginateInMysql(GetConn(), 1, 3, selectSql, sqlExceptSelect, []string{"user_age"}, "%大龙%")
	if err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
}

func TestCreateTypeData(t *testing.T) {
	return
	tableName := "type_test"
	data := make(map[string]interface{})
	data["str_type"] = idkit.CreateUniqueId()
	data["tinyint_type"] = 18
	unixTime, createTime, _ := timekit.GetNowTimeMsAndDate(timekit.DateFormat_YYYY_MM_DD_HH_MM_SS)
	data["bigint_type"] = unixTime
	data["timestamp_type"] = createTime

	result, id, err := SaveInMysql(GetConn(), tableName, data)
	if !result && err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	fmt.Printf("type_test表保存记录返回的主键结果 id=%d\n", id)
}

func TestGetTypeData(t *testing.T) {
	return
	sql := `SELECT id, str_type, tinyint_type, bigint_type, timestamp_type FROM type_test ORDER BY ID DESC LIMIT 1`
	result, err := FindFirstInMysql(GetConn(), sql, nil)
	if err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	fmt.Println(result)
}
