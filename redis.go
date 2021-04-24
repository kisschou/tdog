package tdog

import (
	"context"
	redisImpl "github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

/**
 * Redis模块
 *
 * @Author: Kisschou
 * @Build: 2021-04-24
 */
type redisModel struct {
	engineList map[string]*redisImpl.Client // 引擎列表
	Engine     *redisImpl.Client            // 默认引擎
	db         int                          // 当前使用db
}

/**
 * 初始化Redis模块
 *
 * @return *redisModel
 */
func NewRedis() *redisModel {
	redis := new(redisModel)
	redis.Engine = redis.Change("master")
	return redis
}

/**
 * 引擎切换
 *
 * @param string name 切换的引擎名
 *
 * @return *redis.Client
 */
func (r *redisModel) Change(name string) *redisImpl.Client {
	var host, port, pass string
	var poolSize int = 1
	resultImpls := NewConfig().SetFile("cache").SetPrefix(name+".").GetMulti("host", "port", "pass", "pool_size")
	for k, resultImpl := range resultImpls {
		switch k {
		case "host":
			host = resultImpl.ToString()
			break
		case "port":
			port = resultImpl.ToString()
			break
		case "pass":
			pass = resultImpl.ToString()
			break
		case "pool_size":
			poolSize = resultImpl.ToInt()
			break
		}
	}
	if host == "" || port == "" {
		NewLogger().Error("未找到该标签下的redis配置信息")
		return nil
	}
	engine := redisImpl.NewClient(&redisImpl.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
		PoolSize: poolSize,
	})
	if _, err := engine.Ping(Ctx).Result(); err != nil {
		NewLogger().Error("(" + name + ")redis连接失败:" + err.Error())
		return nil
	}
	// 加入连接池
	if r.engineList == nil {
		r.engineList = make(map[string]*redisImpl.Client, 0)
	}
	r.engineList[name] = engine
	r.db = 0
	r.Engine = engine
	return engine
}

/**
 * 新增其他redis引擎
 *
 * @param string name 引擎名
 * @param string host 连接地址
 * @param string port 端口
 * @param string pass 密码
 * @param int poolSize 连接池大小
 *
 * @return *redis.Client
 */
func (r *redisModel) New(name, host, port, pass string, poolSize int) *redisImpl.Client {
	engine := redisImpl.NewClient(&redisImpl.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
		PoolSize: poolSize,
	})
	if _, err := engine.Ping(Ctx).Result(); err != nil {
		NewLogger().Error("(" + name + ")redis连接失败:" + err.Error())
		return nil
	}
	// 加入连接池
	if r.engineList == nil {
		r.engineList[name] = engine
	}
	r.db = 0
	r.Engine = engine
	return engine
}

/**
 * 库切换
 *
 * @param string name 切换的库
 *
 * @return *redis.Client
 */
func (r *redisModel) Db(index int) *redisImpl.Client {
	r.db = index
	r.Engine = redisImpl.NewClient(&redisImpl.Options{DB: index})
	return r.Engine
}
