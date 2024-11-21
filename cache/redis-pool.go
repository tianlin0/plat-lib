package cache

//import (
//	"fmt"
//	"github.com/gomodule/redigo/redis"
//	"net"
//	"sync"
//)
//
////redis的连接池
//var redisPoolMap = sync.Map{}
//var redisPoolLock sync.Mutex //redis建立连接锁
//var oneRedisLock sync.Mutex
//
//func _getRedis(con *database.Connect) (redis.Conn, error) {
//	if con == nil {
//		con = defaultRedis
//	}
//	redisStr := con.GetConnString()
//
//	if redisStr == "" {
//		return nil, fmt.Errorf("getRedis url error")
//	}
//
//	loggers := logs.DefaultLogger()
//
//	//多个请求并发，这里需要加锁
//	oneRedisLock.Lock()
//	defer oneRedisLock.Unlock()
//
//	dialOptList := make([]redis.DialOption, 0)
//
//	if !utils.IsNumeric(con.Database) {
//		con.Database = "0"
//	}
//	dbName, _ := int(conv.Int64(con.Database))
//
//	dialOptList = append(dialOptList, redis.DialDatabase(dbName))
//	if con.Username != "" {
//		dialOptList = append(dialOptList, redis.DialUsername(con.Username))
//	}
//	if con.Password != "" {
//		dialOptList = append(dialOptList, redis.DialPassword(con.Password))
//	}
//	if con.Extend != nil {
//		if useTLS, ok := con.Extend["useTLS"]; ok {
//			if useTlSBool, ok := useTLS.(bool); ok {
//				dialOptList = append(dialOptList, redis.DialUseTLS(useTlSBool))
//			}
//		}
//	}
//
//	c, err1 := redis.Dial("tcp", net.JoinHostPort(con.Host, con.Port), dialOptList...)
//	if err1 != nil {
//		loggers.Error("DialURL" + err1.Error())
//		c, err1 = redis.DialURL(redisStr, dialOptList...)
//		if err1 != nil {
//			return nil, err1
//		}
//	}
//
//	loggers.Debug("[redis.DialURL end]")
//
//	if defaultRedis == nil {
//		SetRedisCache(con)
//	}
//
//	return c, nil
//}
//
//func getRedisPool(con *database.Connect) (*redis.Pool, error) {
//	if con == nil {
//		con = defaultRedis
//	}
//	redisStr := con.GetConnString()
//	loggers := logs.DefaultLogger()
//	loggers.Debug("getRedisPool redisStr:", redisStr)
//	if redisStr == "" {
//		return nil, fmt.Errorf("getRedis url error")
//	}
//
//	cacheKey := utils.Md5(redisStr)
//
//	dbInstanceTemp, has := redisPoolMap.Load(cacheKey)
//	if has {
//		dbPool, ok := dbInstanceTemp.(*redis.Pool)
//		if ok {
//			return dbPool, nil
//		}
//	}
//
//	//多个请求并发，这里需要加锁
//	redisPoolLock.Lock()
//	defer redisPoolLock.Unlock()
//
//	dbInstanceTemp, has = redisPoolMap.Load(cacheKey)
//	if has {
//		dbPool, ok := dbInstanceTemp.(*redis.Pool)
//		if ok {
//			return dbPool, nil
//		}
//	}
//
//	_, err := _getRedis(con)
//	if err != nil {
//		loggers.Error("_getRedis:", err)
//		return nil, err
//	}
//
//	loggers.Debug("[database getRedisPool]", redisStr, con)
//
//	pool := &redis.Pool{
//		// 最大空闲连接数。
//		MaxIdle: 16,
//		// 当为0时，池中的连接数没有限制。
//		MaxActive: 20,
//		//连接关闭时间 300秒 （300秒不使用自动关闭）
//		IdleTimeout: 300,
//		//连接的redis数据库
//		Dial: func() (redis.Conn, error) {
//			return _getRedis(con)
//		},
//	}
//
//	redisPoolMap.Store(cacheKey, pool)
//
//	return pool, nil
//}
//
//func getOneRedis(con *database.Connect) (redis.Conn, error) {
//	pool, err := getRedisPool(con)
//	if err != nil {
//		logs.DefaultLogger().Error("getOneRedis", err)
//		return nil, err
//	}
//	redisConn := pool.Get()
//	return redisConn, redisConn.Err()
//}
