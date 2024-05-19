如果是多个Insert/Update/Delete操作,要遵循ACID原则,必须用显示事务,以下是使用显示事务的基本步骤:

myDbCon := GetMysqlCon() //获取连接
if myDbCon == nil {
    return false, logiccode.DbConErrorCode()
}
if sql == "" || data == nil {
    return false, logiccode.DbUpdateErrorCode()
}

dbLogger.Debug("%s", sql)

tx, err := myDbCon.Begin() //开启事务
if err != nil {
    dbLogger.Error("%v", err)
    return false, logiccode.DbUpdateErrorCode()
}
defer tx.Commit() //提交事务

stmt, err := tx.Prepare(sql)
if err != nil {
    tx.Rollback() //事务回滚
    dbLogger.Error("%v", err)
    return false, logiccode.DbUpdateErrorCode()
}
defer stmt.Close() //释放连接

result, err := stmt.Exec(data...)
if err != nil {
    tx.Rollback() //事务回滚
    dbLogger.Error("%v", err)
    return false, logiccode.DbUpdateErrorCode()
}

rowsNum, _ := result.RowsAffected()
if rowsNum == 0 {
    return false, logiccode.DbZeroErrorCode()
}

return true, nil