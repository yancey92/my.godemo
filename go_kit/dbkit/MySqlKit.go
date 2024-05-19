package dbkit

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.gumpcome.com/common/go_kit/logiccode"
	"gitlab.gumpcome.com/common/go_kit/logickit"
	"gitlab.gumpcome.com/common/go_kit/strkit"
	"strconv"
	"strings"
	"time"
)

var dbs = make(map[string]*sql.DB)
var ErrNoRows = sql.ErrNoRows

type Page struct {
	PageNumber int         `json:"page_number" desc:"第几页"`
	PageSize   int         `json:"page_size" desc:"每页显示记录数"`
	TotalPage  int         `json:"total_page" desc:"共多少页"`
	TotalRow   int         `json:"total_row" desc:"多少条记录"`
	List       interface{} `json:"list" desc:"分页结果集"`
}

//type PageList struct {
//	PageNumber int                      `json:"page_number" desc:"第几页"`
//	PageSize   int                      `json:"page_size" desc:"每页显示记录数"`
//	TotalPage  int                      `json:"total_page" desc:"共多少页"`
//	TotalRow   int                      `json:"total_row" desc:"多少条记录"`
//	List       []map[string]interface{} `json:"list" desc:"分页结果集"`
//}

func addDbCfg(cfgName string, db *sql.DB) error {
	if cfgName == "" {
		return errors.New("mysql config name is nil!")
	}

	if db == nil {
		return errors.New("mysql conn is nil!")
	}

	if dbs[cfgName] != nil {
		return errors.New("config name is existed!")
	}

	dbs[cfgName] = db
	return nil
}

// @Title 初始化MySQL数据库
// @param userName 	用户名
// @param userPwd 	密码
// @param host 		地址
// @param dbName 	数据库名称
// @param cfgName	数据库连接池配置名称
// @param maxIdle 	最大活跃连接数
// @param maxActive	最大连接数
func InitMysql(userName string, userPwd string, host string, dbName string, cfgName string, maxIdle int, maxActive int) {
	if userName == "" || userPwd == "" || host == "" || dbName == "" {
		panic(fmt.Sprintf("%s mysql connection info is empty!", cfgName))
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local", userName, userPwd, host, dbName))

	if err != nil {
		panic(err.Error())
	}
	if maxIdle <= 0 {
		maxIdle = 10
	}
	if maxActive <= 0 {
		maxActive = 20
	}
	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxActive)
	//设置连接的存活时间,不是指sleep时间,而是从创建连接时算起。
	db.SetConnMaxLifetime(3 * time.Minute)
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	addDbCfg(cfgName, db)
	beego.Info(fmt.Sprintf("%s 数据库初始化成功...", cfgName))
}

// @Title 获取MySQL连接
func GetMysqlCon(cfgName string) (*sql.DB, error) {
	if cfgName == "" {
		return nil, logiccode.DbConfigNameErrorCode()
	}
	return dbs[cfgName], nil
}

// INSERT INTO `user`(name,age,email,gender,height,interests) VALUES (?,?,?,?,?,?)
func CreateMysqlInsertSQL(tableName string, data map[string]interface{}) (string, []interface{}) {
	dataLen := len(data)
	if dataLen <= 0 {
		return "", nil
	}

	params := make([]interface{}, 0)

	//构建INSERT部分的SQL格式
	insertStrBuilder := strkit.StringBuilder{}
	insertStrBuilder.Append("INSERT INTO `").Append(tableName).Append("`(")

	//构建VALUES部分的SQL格式
	valuesStrBuilder := strkit.StringBuilder{}
	valuesStrBuilder.Append(") VALUES (")

	for k, v := range data {
		if len(params) > 0 {
			insertStrBuilder.Append(", ")
			valuesStrBuilder.Append(", ")
		}
		insertStrBuilder.Append("`").Append(k).Append("`")
		valuesStrBuilder.Append("?")
		params = append(params, v)
	}
	valuesStrBuilder.Append(")")

	sql := strkit.StrJoin(insertStrBuilder.ToString(), valuesStrBuilder.ToString())
	return sql, params
}

// UPDATE `user` SET `name` = ? WHERE `id` = ?
func CreateMysqlUpdateByIdSQL(tableName string, data map[string]interface{}) (string, []interface{}) {
	dataLen := len(data)
	if dataLen <= 0 {
		return "", nil
	}

	params := make([]interface{}, 0)

	//构建UPDATE部分的SQL格式
	updateStrBuilder := strkit.StringBuilder{}
	updateStrBuilder.Append("UPDATE `").Append(tableName).Append("` SET ")

	for k, v := range data {
		if k != "id" {
			if len(params) > 0 {
				updateStrBuilder.Append(", ")
			}
			updateStrBuilder.Append("`").Append(k).Append("` = ? ")
			params = append(params, v)
		}
	}
	updateStrBuilder.Append(" WHERE `id` = ?")
	params = append(params, data["id"])

	return updateStrBuilder.ToString(), params
}

// DELETE FROM `user` WHERE `id` = ?
func CreateDeleteMysqlSQL(tableName string, data map[string]interface{}) (string, []interface{}) {
	dataLen := len(data)
	if dataLen <= 0 {
		return "", nil
	}

	params := make([]interface{}, 0)

	//构建DELETE部分的SQL格式
	deleteStrBuilder := strkit.StringBuilder{}
	deleteStrBuilder.Append("DELETE FROM `").Append(tableName).Append("` WHERE ")

	for k, v := range data {
		if len(params) > 0 {
			deleteStrBuilder.Append(" AND ")
		}
		deleteStrBuilder.Append("`").Append(k).Append("` = ?")
		params = append(params, v)
	}

	return deleteStrBuilder.ToString(), params
}

// @Title 保存数据
// @Description 	返回的int64类型的值,只有在表主键定义为"auto increment"情况下,才会有效,其他情况默认返回0
// @param myDbCon	数据库连接
// @param tableName	表名称
// @param data		需要保存的K-V键值对,K:字段名,V:字段值
func SaveInMysql(myDbCon *sql.DB, tableName string, data map[string]interface{}) (bool, int64, error) {
	if myDbCon == nil {
		return false, 0, logiccode.DbConErrorCode()
	}
	if tableName == "" || data == nil {
		return false, 0, logiccode.DbInsertErrorCode()
	}
	sql, params := CreateMysqlInsertSQL(tableName, data)

	beego.Debug(fmt.Sprintf("SQL %s VALS %#v", sql, params))

	result, err := myDbCon.Exec(sql, params...)
	if err != nil {
		beego.Error(fmt.Sprintf("%v", err))
		return false, 0, logiccode.DbInsertErrorCode()
	}

	id, err := result.LastInsertId()

	return true, id, err
}

// @Title 根据主键ID更新数据
// @Description 主键字段名称必须是"id"
// @param myDbCon	数据库连接
// @param tableName	表名称
// @param data		需要保存的K-V键值对,K:字段名,V:字段值
func UpdateByIdInMysql(myDbCon *sql.DB, tableName string, data map[string]interface{}) (bool, error) {
	if data["id"] == nil {
		return false, logiccode.DbUpdateByIdErrorCode()
	}
	if myDbCon == nil {
		return false, logiccode.DbConErrorCode()
	}
	if tableName == "" || data == nil {
		return false, logiccode.DbUpdateErrorCode()
	}
	sql, params := CreateMysqlUpdateByIdSQL(tableName, data)

	beego.Debug(fmt.Sprintf("SQL %s VALS %#v", sql, params))

	result, err := myDbCon.Exec(sql, params...)
	if err != nil {
		beego.Error(fmt.Sprintf("%v", err))
		return false, logiccode.DbUpdateErrorCode()
	}

	rowsNum, _ := result.RowsAffected()
	if rowsNum == 0 {
		return true, logiccode.DbZeroErrorCode()
	}

	return true, nil
}

// @Title 更新数据
// @Description 更新的字段值记录必须与更新SQL语句需要更新的字段顺序一致
// @param myDbCon	数据库连接
// @param sql		更新SQL语句
// @param data		需要更新的字段值记录
func UpdateInMysql(myDbCon *sql.DB, sql string, data ...interface{}) (bool, error) {
	if myDbCon == nil {
		return false, logiccode.DbConErrorCode()
	}
	if sql == "" {
		return false, logiccode.DbUpdateErrorCode()
	}

	beego.Debug(fmt.Sprintf("SQL %s VALS %#v", sql, data))

	result, err := myDbCon.Exec(sql, data...)
	if err != nil {
		beego.Error(fmt.Sprintf("%v", err))
		return false, logiccode.DbUpdateErrorCode()
	}

	rowsNum, _ := result.RowsAffected()
	if rowsNum == 0 {
		return true, logiccode.DbZeroErrorCode()
	}

	return true, nil
}

// @Title 根据主键ID更新数据
// @Description 主键字段名称必须是"id"
// @param myDbCon	数据库连接
// @param tableName	表名称
// @param id		主键ID的字段值
func DeleteByIdInMysql(myDbCon *sql.DB, tableName string, id interface{}) (bool, error) {
	if myDbCon == nil {
		return false, logiccode.DbConErrorCode()
	}
	if tableName == "" || id == nil {
		return false, logiccode.DbDeleteErrorCode()
	}
	sql, params := CreateDeleteMysqlSQL(tableName, map[string]interface{}{"id": id})

	beego.Debug(fmt.Sprintf("SQL %s VALS %#v", sql, params))

	result, err := myDbCon.Exec(sql, params...)
	if err != nil {
		beego.Error(fmt.Sprintf("%v", err))
		return false, logiccode.DbDeleteErrorCode()
	}

	rowsNum, _ := result.RowsAffected()
	if rowsNum == 0 {
		return true, logiccode.DbZeroErrorCode()
	}

	return true, nil
}

// @Title 删除数据
// @Description data保存的参数值必须与删除SQL语句WHERE条件需要的字段顺序一致
// @param myDbCon 	数据库连接
// @param sql		删除SQL语句
// @param data		WHERE条件字段值记录
func DeleteInMysql(myDbCon *sql.DB, sql string, data ...interface{}) (bool, error) {
	if myDbCon == nil {
		return false, logiccode.DbConErrorCode()
	}
	if sql == "" {
		return false, logiccode.DbDeleteErrorCode()
	}

	beego.Debug(fmt.Sprintf("SQL %s VALS %#v", sql, data))

	result, err := myDbCon.Exec(sql, data...)
	if err != nil {
		beego.Error(fmt.Sprintf("%v", err))
		return false, logiccode.DbDeleteErrorCode()
	}

	rowsNum, _ := result.RowsAffected()
	if rowsNum == 0 {
		return true, logiccode.DbZeroErrorCode()
	}

	return true, nil
}

// @Title 查询数据
// @Description data保存的参数值必须与查询SQL语句WHERE条件需要的字段顺序一致
// @param myDbCon 	数据库连接
// @param sql		查询SQL语句
// @param intItems	需要转换成Int类型的字段集合
// @param data		WHERE条件字段值记录
func FindInMysql(myDbCon *sql.DB, querySql string, intItems []string, data ...interface{}) ([]map[string]interface{}, error) {
	if myDbCon == nil {
		return nil, logiccode.DbConErrorCode()
	}
	if querySql == "" {
		return nil, logiccode.DbDeleteErrorCode()
	}
	intItemsMap := make(map[string]int)
	if intItems != nil && len(intItems) > 0 {
		for idx, item := range intItems {
			intItemsMap[item] = idx
		}
	}

	beego.Debug(fmt.Sprintf("SQL %s VALS %#v", querySql, data))

	rows, err := myDbCon.Query(querySql, data...)

	if err == sql.ErrNoRows {
		//没有查到结果
		return nil, nil
	}
	if err != nil {
		beego.Error(fmt.Sprintf("%v", err))
		return nil, logiccode.DbQueryErrorCode()
	}

	defer rows.Close()

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var records = make([]map[string]interface{}, 0, 10)
	var itemVal string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]interface{})
		for i, col := range values {
			itemVal = ""
			switch col.(type) {
			case []uint8:
				itemVal = string(col.([]uint8))
			case int64:
				itemVal = fmt.Sprint(col.(int64))
			case float32:
				itemVal = fmt.Sprint(col.(float32))
			}

			if _, ok := intItemsMap[columns[i]]; ok {
				if strings.TrimSpace(itemVal) != "" {
					record[columns[i]], err = strkit.StrToInt(itemVal)
					if err != nil {
						beego.Error(fmt.Sprintf("item to int error :%#v", columns[i]))
						return nil, logiccode.DbItemToIntErrorCode()
					}
				} else {
					record[columns[i]] = 0
				}
			} else {
				record[columns[i]] = itemVal
			}
		}
		if len(record) > 0 {
			records = append(records, record)
		}
	}
	if len(records) == 0 {
		records = nil
	}
	return records, err
}

// @Title 查询单条数据
// @Description data保存的参数值必须与查询SQL语句WHERE条件需要的字段顺序一致,如果查询SQL影响的行数多与1行,必须追加 LIMIT 1 条件
// @param myDbCon 	数据库连接
// @param sql		查询SQL语句
// @param intItems	需要转换成Int类型的字段集合
// @param data		WHERE条件字段值记录
func FindFirstInMysql(myDbCon *sql.DB, querySql string, intItems []string, data ...interface{}) (map[string]interface{}, error) {
	if myDbCon == nil {
		return nil, logiccode.DbConErrorCode()
	}
	if querySql == "" {
		return nil, logiccode.DbDeleteErrorCode()
	}
	intItemsMap := make(map[string]int)
	if intItems != nil && len(intItems) > 0 {
		for idx, item := range intItems {
			intItemsMap[item] = idx
		}
	}

	beego.Debug(fmt.Sprintf("SQL %s VALS %#v", querySql, data))

	rows, err := myDbCon.Query(querySql, data...)

	if err == sql.ErrNoRows {
		//没有查到结果
		return nil, nil
	}
	if err != nil {
		beego.Error(fmt.Sprintf("%v", err))
		return nil, logiccode.DbQueryErrorCode()
	}

	defer rows.Close()

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var records = make([]map[string]interface{}, 0, 10)
	var itemVal string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]interface{})
		for i, col := range values {
			itemVal = ""
			switch col.(type) {
			case []uint8:
				itemVal = string(col.([]uint8))
			case int64:
				itemVal = fmt.Sprint(col.(int64))
			}

			if _, ok := intItemsMap[columns[i]]; ok {
				if strings.TrimSpace(itemVal) != "" {
					record[columns[i]], err = strkit.StrToInt(itemVal)
					if err != nil {
						beego.Error(fmt.Sprintf("item to int error :%#v", columns[i]))
						return nil, logiccode.DbItemToIntErrorCode()
					}
				} else {
					record[columns[i]] = 0
				}
			} else {
				record[columns[i]] = itemVal
			}
		}
		if len(record) > 0 {
			records = append(records, record)
		}
	}
	if len(records) == 0 {
		return nil, nil
	}

	return records[0], nil
}

// @Tile 分页查询
// @param myDbCon		数据库连接
// @param pageNumber		第几页,最小值为1
// @param pageSize		每页显示几条记录,最多100条
// @param selectSql		查询SQL
// @param sqlExceptSelect	查询SQL条件
// @param intItems	需要转换成Int类型的字段集合
// @param data		WHERE条件字段值记录
func PaginateInMysql(myDbCon *sql.DB, pageNumber int, pageSize int, selectSql string, sqlExceptSelect string, intItems []string, data ...interface{}) (Page, error) {
	if myDbCon == nil {
		return Page{}, logiccode.DbConErrorCode()
	}

	if pageNumber < 0 || pageSize <= 0 || !logickit.VerificationPageSize(pageSize) {
		return Page{}, logiccode.DbPageOutErrorCode()
	}

	//统计记录总数
	totalRowSqlBuilder := strkit.StringBuilder{}
	totalRowSqlBuilder.Append("SELECT COUNT(*) AS count ").Append(sqlExceptSelect)
	totalRowResult, err := FindFirstInMysql(myDbCon, totalRowSqlBuilder.ToString(), []string{"count"}, data...)

	if err != nil {
		return Page{}, err
	}

	totalRow, ok := totalRowResult["count"].(int)
	if !ok {
		return Page{}, logiccode.DbPageCountToIntCode()
	}

	//计算共多少页记录
	totalPage := totalRow / pageSize
	if totalRow%pageSize != 0 {
		totalPage++
	}

	//查询的页码超出总页数
	if pageNumber > totalPage {
		return Page{PageNumber: pageNumber, PageSize: pageSize, TotalPage: totalPage, TotalRow: totalRow}, nil
	}

	offset := pageSize * (pageNumber - 1)

	pageSqlBuilder := strkit.StringBuilder{}
	pageSqlBuilder.Append(selectSql).Append(" ").Append(sqlExceptSelect).Append(" LIMIT ").Append(strconv.Itoa(offset)).Append(", ").Append(strconv.Itoa(pageSize))
	pageResult, err := FindInMysql(myDbCon, pageSqlBuilder.ToString(), intItems, data...)
	if err != nil {
		return Page{}, err
	}

	return Page{PageNumber: pageNumber, PageSize: pageSize, TotalPage: totalPage, TotalRow: totalRow, List: pageResult}, nil
}

//func PageListMysql(myDbCon *sql.DB, pageNumber int, pageSize int, selectSql string, sqlExceptSelect string, intItems []string, data ...interface{}) (*PageList, error) {
//	result := new(PageList)
//	result.PageNumber = pageNumber
//	result.PageSize = pageSize
//	if myDbCon == nil {
//		return nil, logiccode.DbConErrorCode()
//	}
//
//	if pageNumber < 0 || pageSize <= 0 || pageSize > 100 {
//		return nil, logiccode.DbPageOutErrorCode()
//	}
//
//	//统计记录总数
//	totalRowSqlBuilder := strkit.StringBuilder{}
//	totalRowSqlBuilder.Append("SELECT COUNT(*) AS count ").Append(sqlExceptSelect)
//	totalRowResult, err := FindFirstInMysql(myDbCon, totalRowSqlBuilder.ToString(), []string{"count"}, data...)
//
//	if err != nil {
//		return nil, err
//	}
//
//	totalRow, ok := totalRowResult["count"].(int)
//	if !ok {
//		return nil, logiccode.DbPageCountToIntCode()
//	}
//
//	//计算共多少页记录
//	totalPage := totalRow / pageSize
//	if totalRow%pageSize != 0 {
//		totalPage++
//	}
//
//	result.TotalRow = totalRow
//	result.TotalPage = totalPage
//
//	//查询的页码超出总页数
//	if pageNumber > totalPage {
//		return result, nil
//	}
//
//	offset := pageSize * (pageNumber - 1)
//
//	pageSqlBuilder := strkit.StringBuilder{}
//	pageSqlBuilder.Append(selectSql).Append(" ").Append(sqlExceptSelect).Append(" LIMIT ").Append(strconv.Itoa(offset)).Append(", ").Append(strconv.Itoa(pageSize))
//	pageResult, err := FindInMysql(myDbCon, pageSqlBuilder.ToString(), intItems, data...)
//	if err != nil {
//		return nil, err
//	}
//	result.List = pageResult
//
//	return result, nil
//}
