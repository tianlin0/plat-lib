package database

import (
	"fmt"
	"github.com/tianlin0/plat-lib/conn"
	"github.com/tianlin0/plat-lib/crontab"
	"github.com/tianlin0/plat-lib/encode"
	"github.com/tianlin0/plat-lib/goroutines"
	"github.com/tianlin0/plat-lib/logs"
	"github.com/tianlin0/plat-lib/utils"
	"runtime"
	"sync"
	"time"
	"xorm.io/xorm"
	//需要引入默认的mysql数据驱动
	_ "github.com/go-sql-driver/mysql"
)

var dbInstanceMap = new(sync.Map)
var dbConnectMap = sync.Map{}

var dbMux sync.Mutex
var once sync.Once

// 异步刷新mysql连接
func syncMysqlPing() {
	goroutines.GoSyncHandler(func(params ...interface{}) {
		dbInstanceMap.Range(func(key, value interface{}) bool {
			dbInstance, ok := value.(*xorm.Engine)
			if !ok {
				//判断是否是mysql类型，如果是，则删除
				if connTemp, ok2 := dbConnectMap.Load(key); ok2 {
					if connTemp2, ok3 := connTemp.(*conn.Connect); ok3 {
						if connTemp2.Driver == conn.DriverMysql {
							dbInstanceMap.Delete(key)
						}
					}
				}
				return true
			}

			//检测连接是否存活中，如果是通的话，则直接返回，失败则重新连接
			//每次检测会影响效率，增加一个概率功能
			ranInt := utils.Random(0, 100)
			if ranInt > 20 { //20%的概率会检测
				return true
			}
			err := dbInstance.Ping()
			if err == nil {
				return true
			}

			logs.DefaultLogger().Error("dbInstance.Ping:", err.Error())
			//重新建立连接
			if connTemp, ok2 := dbConnectMap.Load(key); ok2 {
				if connTemp2, ok3 := connTemp.(*conn.Connect); ok3 {
					if connTemp2.Driver == conn.DriverMysql {
						dbInstanceMap.Delete(key) //先删除才能创建新的
						engines := GetXormEngine(connTemp2)
						if engines != nil {
							err = dbInstance.Close() //存在新的话，则原来的就需要断开连接
							if err != nil {
								logs.DefaultLogger().Error("dbInstance.Close error:", err.Error())
							}
						} else {
							logs.DefaultLogger().Error("dbInstanceMap/GetXormEngine nil:", connTemp2)
						}
					}
				}
			}

			return true
		})
	}, nil)
}

func getMysqlCacheKey(con *conn.Connect) string {
	keyStr, err := encode.Serialize(con)
	if err != nil {
		conStr := con.GetConnString()
		keyStr = fmt.Sprintf("%s://%s", con.Driver, conStr)
	}
	return encode.Md5(keyStr)
}

// GetXormEngine 获取数据库单例对象
func GetXormEngine(con *conn.Connect) *xorm.Engine {
	if con == nil {
		return nil
	}

	cacheKey := getMysqlCacheKey(con)
	dbInstanceTemp, has := dbInstanceMap.Load(cacheKey)
	if has {
		dbInstance, ok := dbInstanceTemp.(*xorm.Engine)
		if ok {
			return dbInstance
		}
	}

	dbMux.Lock()
	defer dbMux.Unlock()

	//带检查锁的单例模式
	dbInstanceTemp, has = dbInstanceMap.Load(cacheKey)
	if has {
		dbInstance, ok := dbInstanceTemp.(*xorm.Engine)
		if ok {
			return dbInstance
		}
	}

	if con.Driver == conn.DriverMysql {
		connStr := con.GetConnString()
		engine, err := xorm.NewEngine(con.Driver, connStr)
		err2 := engine.Ping()
		logs.DefaultLogger().Debug("mysql.GetXormEngine:", con, engine, err, connStr, err2)
		if err != nil {
			logs.DefaultLogger().Error("dbHelper.GetXormEngine error=", err)
			return nil
		}
		engine.SetConnMaxLifetime(0) //设置 packets.go:123: closing bad idle connection: EOF
		// 设置连接池的最大空闲连接数
		engine.SetMaxIdleConns(10) //默认为2
		// 设置连接池的最大打开连接数
		//engine.SetMaxOpenConns(100) //默认无限

		dbInstanceMap.Store(cacheKey, engine)
		dbConnectMap.Store(cacheKey, con)

		//定时任务执行,检查mysql连接是否正常
		once.Do(func() {
			//当dbInstanceMap销毁时，则需要断开连接
			runtime.SetFinalizer(dbInstanceMap, func(dbInstanceMapPrt *sync.Map) {
				if dbInstanceMapPrt == nil {
					return
				}
				dbInstanceMapPrt.Range(func(key, value interface{}) bool {
					dbInstance, ok := value.(*xorm.Engine)
					if ok {
						err = dbInstance.Close()
						if err != nil {
							logs.DefaultLogger().Error("GetXormEngine SetFinalizer close errr:", err)
						}
						dbInstanceMapPrt.Delete(key)
					}
					return true
				})
			})

			crontab.StartCrontabJobs(60*time.Second, map[string]func(){
				"*/7 * * * *": func() {
					//每7分钟执行数据库连接检测
					syncMysqlPing()
				},
			})
		})
		return engine
	}

	return nil
}
