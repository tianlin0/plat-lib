package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tianlin0/plat-lib/cond"
	"github.com/tianlin0/plat-lib/conn"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/logs"
	"github.com/tianlin0/plat-lib/utils"
	"reflect"
	"strings"
	"sync"
	"xorm.io/core"
	"xorm.io/xorm"
)

// Dao 访问数据库的对象
type Dao struct {
	once    sync.Once
	connect *conn.Connect
	ctx     context.Context
	Engine  *xorm.Engine
	//如果有事务，则将使用
	daoSession *xorm.Session
}

var (
	explainSql = false //执行分析索引命中的情况
)

// TransCallback 事务回调函数
type TransCallback func(*xorm.Session) error

// SetTableSuffix 设置表的后缀，table_1, table_2等
func (d *Dao) SetTableSuffix(suffix string) {
	tbMapper := core.NewSuffixMapper(core.GonicMapper{}, suffix)
	d.Engine.SetTableMapper(tbMapper)
}

// SetLogger 设置日志
func (d *Dao) SetLogger(loggerOld interface{}) {
	logger := setXormLogger(loggerOld)
	//一个链接只需要执行一次
	if logger != nil {
		d.Engine.SetLogger(logger)
		d.Engine.ShowSQL(true)
	} else {
		d.Engine.ShowSQL(false)
	}
}

// SetExplainSql 设置是否需要调试
func SetExplainSql(explain bool) {
	explainSql = explain
}

// setLogger 强制设置logger
func (d *Dao) setLogger() {
	d.SetLogger(logs.CtxLogger(d.ctx))
}

// initDB 初始化连接，内部
func (d *Dao) initDB(ctx context.Context, co *conn.Connect) *Dao {
	engine := GetXormEngine(co)
	if engine == nil {
		logs.CtxLogger(ctx).Error("initDB/GetXormEngine nil:", co)
		return nil
	}
	d.Engine = engine
	d.connect = co
	d.ctx = ctx

	d.once.Do(func() {
		//默认打印
		d.SetLogger(logs.DefaultLogger())
	})
	return d
}

func checkInitDB(ctx context.Context, child interface{}, co *conn.Connect, isPanic bool) *Dao {
	obj, ok := child.(interface {
		initDB(ctx context.Context, con *conn.Connect) *Dao
	})
	if ok {
		ret := obj.initDB(ctx, co)
		if ret == nil {
			logs.CtxLogger(ctx).Error("InitDatabase error:", co)
			if isPanic {
				panic(interface{}(co))
			}
			return nil
		}
		//进行特殊格式打印？
		ret.setLogger()
		return ret
	}
	return nil
}

// BuildXormEngine 对外调用继承Dao的对象，进行数据库连接初始化
func BuildXormEngine(ctx context.Context, child interface{}, con *conn.Connect, isPanic ...bool) (*xorm.Engine, error) {
	isPanicBool := false
	if len(isPanic) >= 1 {
		isPanicBool = isPanic[0]
	}

	ret := checkInitDB(ctx, child, con, isPanicBool)
	if ret == nil {
		logs.CtxLogger(ctx).Error("child error:", con)
		return nil, fmt.Errorf("InitDatabase child error")
	}
	return ret.Engine, nil
}

// Insert 新增
func (d *Dao) Insert(info ...interface{}) (int64, error) {
	if d.daoSession != nil {
		return d.daoSession.Insert(info...)
	}

	num, err := d.Engine.Insert(info...)
	if err != nil {
		return -1, err
	}
	return num, nil
}

// FlagDelete 逻辑删除
func (d *Dao) FlagDelete(id int64, info interface{}) (int64, error) {
	if d.daoSession != nil {
		return d.daoSession.ID(id).Delete(info)
	}
	return d.Engine.ID(id).Delete(info)
}

// Delete 删除
func (d *Dao) Delete(id interface{}, info interface{}) (int64, error) {
	if d.daoSession != nil {
		return d.daoSession.ID(id).Unscoped().Delete(info)
	}
	return d.Engine.ID(id).Unscoped().Delete(info)
}

// Update 更新
func (d *Dao) Update(id interface{}, info interface{}, columns ...string) (int64, error) {
	if len(columns) == 0 {
		//update 默认值为0的情况不会更新的问题
	}
	if d.daoSession != nil {
		if len(columns) > 0 {
			num, err := d.daoSession.ID(id).Cols(columns...).Update(info)
			if err != nil {

			}
			return num, err
		}
		return d.daoSession.ID(id).Update(info)
	}
	if len(columns) > 0 {
		return d.Engine.ID(id).Cols(columns...).Update(info)
	}
	return d.Engine.ID(id).Update(info)
}

// Get 通过主键查询单个
func (d *Dao) Get(id interface{}, info interface{}) (bool, error) {
	ok, err := d.Engine.ID(id).Get(info)
	if ok {
		if err == nil {
			return true, nil
		}
		return true, err
	}
	return false, err
}

// UpdateWhere 条件更新
func (d *Dao) UpdateWhere(whereStr string, argList []interface{}, info interface{}, columns ...string) (int64, error) {
	if d.daoSession != nil {
		if len(columns) > 0 {
			num, err := d.daoSession.Where(whereStr, argList...).Cols(columns...).Update(info)
			if err != nil {

			}
			return num, err
		}
		return d.daoSession.Where(whereStr, argList...).Update(info)
	}
	if len(columns) > 0 {
		return d.Engine.Where(whereStr, argList...).Cols(columns...).Update(info)
	}
	return d.Engine.Where(whereStr, argList...).Update(info)
}

// DeleteWhere 条件删除
func (d *Dao) DeleteWhere(whereStr string, argList []interface{}, info interface{}) (int64, error) {
	if d.daoSession != nil {
		return d.daoSession.Where(whereStr, argList...).Unscoped().Delete(info)
	}
	return d.Engine.Where(whereStr, argList...).Unscoped().Delete(info)
}

// GetWhere 通过where查询单个
func (d *Dao) GetWhere(whereStr string, argList []interface{}, info interface{}) (bool, error) {
	ok, err := d.Engine.Where(whereStr, argList...).Get(info)
	if ok {
		return true, err
	}
	return false, err
}

// TransAction 事务
func (d *Dao) TransAction(callback TransCallback) error {
	session := d.Engine.NewSession()
	defer func(session *xorm.Session) {
		_ = session.Close()
	}(session)

	if err := session.Begin(); err != nil {
		return fmt.Errorf("fail to session begin：" + err.Error())
	}

	d.daoSession = session
	err := callback(session)
	d.daoSession = nil
	if err != nil {
		_ = session.Rollback()
		return err
	}
	return session.Commit()
}

// GetListByMap 通过对象查询列表
func (d *Dao) GetListByMap(info map[string]interface{}, bean interface{}) ([]map[string]string, error) {
	tableInfo, err := d.Engine.TableInfo(bean)
	if err != nil {
		return nil, err
	}
	newInfo := make(map[string]interface{})
	for name, val := range info {
		isFind := false
		for _, oneColumn := range tableInfo.Columns() {
			if oneColumn.Name == name {
				isFind = true
				break
			}
		}
		if isFind {
			newInfo[name] = val
		}
	}
	tempStatement := d.Engine.Table(bean)
	for key, val := range newInfo {
		compare := "="
		if reflect.TypeOf(val).Kind() == reflect.Slice {
			list := make([]interface{}, 0)
			s := reflect.ValueOf(val)
			if s.Len() == 2 {
				tempCompare := conv.String(s.Index(0))
				compareList := []string{"like", "=", ">=", ">", "<=", "<"}
				isFind := false
				for _, one := range compareList {
					if one == tempCompare {
						isFind = true
						break
					}
				}
				if isFind {
					tempStatement = tempStatement.Where("`"+key+"` "+tempCompare+" ?", val)
					continue
				}
			}
			//需要去重处理
			onlyArray := make([]string, 0)
			for i := 0; i < s.Len(); i++ {
				ele := s.Index(i).Interface()
				tempOne := conv.String(ele)
				if ret, _ := cond.Contains(onlyArray, tempOne); !ret {
					onlyArray = utils.AppendUnique(onlyArray, tempOne)
					list = append(list, ele)
				}
			}
			if len(list) > 0 {
				tempStatement = tempStatement.In(key, list...)
			}
			continue
		}
		tempStatement = tempStatement.Where("`"+key+"`"+compare+"?", val)
	}
	retMap, err := tempStatement.QueryString()
	if err != nil {
		return nil, err
	}
	if retMap != nil {
		return retMap, nil
	}
	return make([]map[string]string, 0), nil
}

func (d *Dao) explainSqlHandle(sqlOrArgs ...interface{}) {
	if len(sqlOrArgs) == 0 {
		return
	}
	sType := getStatementType(conv.String(sqlOrArgs[0]))
	if sType != "SELECT" {
		//不是查询语句不进行分析
		return
	}

	oldSql := sqlOrArgs[0]

	sqlOrArgs[0] = fmt.Sprintf("%s %s", "EXPLAIN", sqlOrArgs[0])
	retList, err := d.Engine.Query(sqlOrArgs...)
	if err != nil {
		logs.DefaultLogger().Error(sqlOrArgs[0], err)
		return
	}
	if len(retList) == 0 {
		return
	}
	for _, one := range retList {
		possibleKeys, ok1 := one["possible_keys"]
		keys, ok2 := one["key"]
		table, ok3 := one["table"]
		if ok1 && ok2 {
			useKey := false //是否使用了索引

			keyStr := string(keys)
			if keyStr != "" {
				possibleKeyList := strings.Split(string(possibleKeys), ",")
				keyList := strings.Split(keyStr, ",")
				if len(keyList) > 0 {
					for _, oneKey := range keyList {
						if oneKey == "" {
							continue
						}
						if ok, _ := cond.Contains(possibleKeyList, oneKey); ok {
							useKey = true
							break
						}
					}
				}
			}

			if !useKey {
				tableName := ""
				if ok3 {
					tableName = string(table)
				}
				logs.DefaultLogger().Error("has no select index:", oldSql, "|",
					tableName, "|", string(possibleKeys), "|", keyStr)
			}
		}
	}

}

func getStatementType(sql string) string {
	if len(sql) == 0 {
		return ""
	}

	// 去除首尾空格并转换为大写
	sql = strings.TrimSpace(strings.ToUpper(sql))

	if strings.HasPrefix(sql, "SELECT") {
		return "SELECT"
	} else if strings.HasPrefix(sql, "INSERT") {
		return "INSERT"
	} else if strings.HasPrefix(sql, "UPDATE") {
		return "UPDATE"
	} else if strings.HasPrefix(sql, "DELETE") {
		return "DELETE"
	}

	return ""
}

// SqlQuery sql查询
func (d *Dao) SqlQuery(sqlStr string, args ...interface{}) ([]map[string]string, error) {
	queryParam := make([]interface{}, 0)
	queryParam = append(queryParam, sqlStr)
	if args != nil && len(args) > 0 {
		queryParam = append(queryParam, args...)
	}
	retList, err := d.Engine.Query(queryParam...)
	if err != nil {
		logs.DefaultLogger().Error("SqlQuery Error:", err, sqlStr, d.Engine)
		return nil, err
	}
	if explainSql {
		d.explainSqlHandle(queryParam...)
	}

	if retList == nil {
		return []map[string]string{}, nil
	}
	retData := make([]map[string]string, len(retList))
	for i, one := range retList {
		oneTemp := make(map[string]string, 0)
		for key, val := range one {
			oneTemp[key] = string(val)
		}
		retData[i] = oneTemp
	}
	return retData, nil
}

// SqlExec sql更新
func (d *Dao) SqlExec(sqlStr string, args ...interface{}) (int64, error) {
	queryParam := make([]interface{}, 0)
	queryParam = append(queryParam, sqlStr)
	if args != nil && len(args) > 0 {
		queryParam = append(queryParam, args...)
	}
	var execResult sql.Result
	var err error
	if d.daoSession != nil {
		execResult, err = d.daoSession.Exec(queryParam...)
	} else {
		execResult, err = d.Engine.Exec(queryParam...)
	}

	if err != nil {
		return 0, err
	}
	num, err := execResult.LastInsertId()
	if err == nil && num > 0 {
		return num, nil
	}

	num, err = execResult.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func getColumnLikeSql(oldValue string, replaceList []string, eacapeList []string) (retValLike string, retEscape string, retSuccess bool) {
	isFind := false
	for _, one := range replaceList {
		if in := strings.IndexAny(oldValue, one); in >= 0 {
			isFind = true
		}
	}

	if !isFind {
		return oldValue, "", true
	}

	oneEscapeStr := ""
	for _, one := range eacapeList {
		if in := strings.IndexAny(oldValue, one); in < 0 {
			//不存在，则可作为转义符
			oneEscapeStr = one
			break
		}
	}

	if oneEscapeStr == "" {
		return oldValue, "", false
	}

	for _, one := range replaceList {
		oldValue = strings.ReplaceAll(oldValue, one, oneEscapeStr+one)
	}

	return oldValue, oneEscapeStr, true
}

// GetColumnLikeSql 获取列名转义sql
func GetColumnLikeSql(oldValue string) (retValLike string, retParam string) {
	replaceList := []string{"%", "_"}
	escapeList := []string{"/", "&", "#", "@", "^", "$", "!"}

	newValue, escape, retTrue := getColumnLikeSql(oldValue, replaceList, escapeList)
	if retTrue {
		if escape != "" {
			return "? escape '" + escape + "'", newValue
		}
	}

	return "?", newValue
}
