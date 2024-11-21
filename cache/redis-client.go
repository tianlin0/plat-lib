package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/tianlin0/plat-lib/cond"
	"github.com/tianlin0/plat-lib/conn"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/logs"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var defaultRedis *conn.Connect

var redisMap = cmap.New()
var onceError sync.Once

var (
	poolMaxSize = 100
	poolMinSize = 10

	poolMinIdleConns = 30            //连接池中最小的空闲连接数，可以通过此属性提供更快的连接分配，默认为0
	poolMaxConnAge   = 3 * time.Hour //Redis 连接的最大寿命，在连接池中的连接达到最大寿命时，客户端会将连接归还到连接池中，
	// 从而避免连接长时间占用资源。默认为不限制连接寿命
	poolPoolTimeout time.Duration = 0 //当连接池中所有连接均被占用时，客户端调用连接池中连接的 Get() 方法会等待的最长时间。
	// 默认值为 ReadTimeout 加上1秒
	poolIdleTimeout = 5 * time.Minute //Redis 连接在空闲状态下的最长存活时间，超过该时间的连接将被关闭。如果指定的值小于服务器上
	// 的超时时间，则客户端在检查连接空闲时会关闭连接，以防止服务器出现连接超时。默认为5分钟。将其设为-1可以禁用连接空闲超时检查
	poolIdleCheckFrequency = time.Minute //空闲连接检查频率。默认为1分钟。将其设为-1可以禁用连接空闲超时检查器，但是仍然会
)

// SetDefaultRedis 切换默认的redis连接
func SetDefaultRedis(con *conn.Connect) {
	if con != nil {
		defaultRedis = con
	}
}

func getOneRedis(con *conn.Connect) (*redis.Client, error) {
	if con == nil {
		con = defaultRedis
	}
	if con == nil {
		return nil, fmt.Errorf("redis conn nil")
	}

	con.Driver = conn.DriverRedis
	redisStr, err := con.GetConnString()

	if redisStr == "" || err != nil {
		return nil, fmt.Errorf("getRedis url error")
	}
	oneCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var closeOldPool *redis.Client

	// 这里需要进行redis初始化
	if data, ok := redisMap.Get(redisStr); ok {
		if oldPool, ok := data.(*redis.Client); ok {
			if !cond.IsNil(oldPool) {
				_, err = oldPool.Ping(oneCtx).Result()
				if err == nil {
					return oldPool, nil
				}

				logs.DefaultLogger().Error("_getRedis:", err)

				closeOldPool = oldPool
			}
		} else {
			redisMap.Remove(redisStr)
		}
	}

	dialOpt := getRedisOption(con)
	pool := redis.NewClient(dialOpt)
	_, err = pool.Ping(oneCtx).Result()
	if err != nil {
		_ = pool.Close()
		return nil, err
	}

	//新建以后，需要回收老的
	if closeOldPool != nil {
		defer func(oldPool *redis.Client) {
			_ = oldPool.Close()
		}(closeOldPool)
	}

	redisMap.Set(redisStr, pool)

	if defaultRedis == nil {
		SetDefaultRedis(con)
	}

	return pool, nil
}

func getRedisOption(con *conn.Connect) *redis.Options {
	dialOpt := &redis.Options{}
	{ //db
		if !cond.IsNumeric(con.Database) {
			con.Database = "0"
		}
		dbName, err := strconv.Atoi(con.Database)
		if err != nil {
			dbName = 0
		}
		dialOpt.DB = dbName
	}

	if con.Username != "" {
		dialOpt.Username = con.Username
	}

	if con.Password != "" {
		dialOpt.Password = con.Password
	}

	if con.Extend != nil {
		if useTLS, ok := con.Extend["useTLS"]; ok {
			if useTlSBool, ok := useTLS.(bool); ok {
				if useTlSBool {
					tlsConfig := &tls.Config{InsecureSkipVerify: true}
					if tlsConfig.ServerName == "" {
						tlsConfig.ServerName = con.Host
					}
					dialOpt.TLSConfig = tlsConfig
				}
			}
		}
	}
	dialOpt.Addr = net.JoinHostPort(con.Host, con.Port)
	if con.Protocol != "" {
		dialOpt.Network = con.Protocol
	}

	{ // 连接池的配置
		dialOpt.PoolFIFO = true                 //Redis 连接池是否使用 FIFO 先进先出的连接池类型，默认为 true
		dialOpt.PoolSize = getPoolSize(con)     //连接池中最多能同时存放的 Redis 连接数，即最大连接数
		dialOpt.MinIdleConns = poolMinIdleConns //连接池中最小的空闲连接数，可以通过此属性提供更快的连接分配，默认为0
		dialOpt.MaxConnAge = poolMaxConnAge     //Redis 连接的最大寿命，在连接池中的连接达到最大寿命时，客户端会将连接归还到连接池中，
		// 从而避免连接长时间占用资源。默认为不限制连接寿命
		dialOpt.PoolTimeout = poolPoolTimeout //当连接池中所有连接均被占用时，客户端调用连接池中连接的 Get() 方法会等待的最长时间。
		// 默认值为 ReadTimeout 加上1秒
		dialOpt.IdleTimeout = poolIdleTimeout //Redis 连接在空闲状态下的最长存活时间，超过该时间的连接将被关闭。如果指定的值小于服务器上
		// 的超时时间，则客户端在检查连接空闲时会关闭连接，以防止服务器出现连接超时。默认为5分钟。将其设为-1可以禁用连接空闲超时检查
		dialOpt.IdleCheckFrequency = poolIdleCheckFrequency //空闲连接检查频率。默认为1分钟。将其设为-1可以禁用连接空闲超时检查器，但是仍然会
		// 根据 IdleTimeout 的值关闭空闲连接。
	}
	return dialOpt
}

func getPoolSize(con *conn.Connect) int {
	poolSize := 0
	if con.Extend != nil {
		if pSize, ok := con.Extend["poolSize"]; ok {
			pSizeInt, _ := conv.Int64(pSize)
			if pSizeInt > 0 {
				poolSize = int(pSizeInt)
			}
		}
	}
	if poolSize <= 0 {
		poolSize = runtime.GOMAXPROCS(0)
	}
	if poolSize < poolMinSize {
		poolSize = poolMinSize
	}
	if poolSize > poolMaxSize {
		poolSize = poolMaxSize
	}
	return poolSize
}

// GetRedisClient 获取redis客户端
func GetRedisClient(con *conn.Connect) (*redis.Client, error) {
	loggers := logs.DefaultLogger()

	cli, err := getOneRedis(con)
	if err != nil {
		if con != nil {
			// 如果未设置redis，则提示
			loggers.Error("[redis-client] error:", con, err.Error())
		} else {
			// 没有设置，全局只提醒一次
			onceError.Do(func() {
				loggers.Warn("[redis-client] no set empty:", err.Error())
			})
		}
		return nil, err
	}
	return cli, nil
}
