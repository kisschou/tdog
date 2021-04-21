package tdog

import (
	RedisModel "github.com/go-redis/redis/v7"
)

type Redis struct {
	Engine *RedisModel.Client
}

func (redis *Redis) NewEngine() {
	ConfigTdog := NewConfig()
	host := ConfigTdog.Get("cache.master.host").ToString()
	port := ConfigTdog.Get("cache.master.port").ToString()
	pass := ConfigTdog.Get("cache.master.pass").ToString()

	redis.Engine = RedisModel.NewClient(&RedisModel.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
	})
}

func (redis *Redis) Ping() bool {
	_, err := redis.Engine.Ping().Result()
	if err != nil {
		NewLogger().Error(err.Error())
		return false
	}
	return true
}

func (redis *Redis) Select(index int) {
	redis.Engine = RedisModel.NewClient(&RedisModel.Options{DB: index})
}
