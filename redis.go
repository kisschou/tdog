package tdog

import (
	RedisImpl "github.com/go-redis/redis/v7"
)

type redisModel struct {
	engineList map[string]*RedisImpl.Client // 引擎列表
	Engine     *RedisImpl.Client            // 默认引擎
}

func NewEngine() *redisModel {

	redis := new(redisModel)
	return redis
}

func (redis *redisModel) Change(name string) *RedisImpl.Client {
	var host, port, pass string
	resultImpls := NewConfig().SetFile("cache").SetPrefix(name+".").GetMulti("host", "port", "pass")
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
		}
	}
	if host == "" || port == "" {
		NewLogger().Error("未找到该标签下的redis配置信息")
		return nil
	}
	return RedisImpl.NewClient(&RedisImpl.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
	})
}

func (redis *redisModel) Db(index int) *RedisImpl.Client {
	redis.Engine = RedisImpl.NewClient(&RedisImpl.Options{DB: index})
	return redis.Engine
}
