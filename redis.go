// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	"context"
	redisImpl "github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

/**
 * The module for redis handler.
 *
 * @Author: Kisschou
 * @Build: 2021-04-24
 */
type redisModel struct {
	engineList map[string]*redisImpl.Client // Engine Pool
	Engine     *redisImpl.Client            // Current Engine
	db         int                          // Current Db
}

// NewRedis init redis module
// auto-load master configuration and make it actived
// so...must write master configuration of redis into configuration.
func NewRedis() *redisModel {
	redis := new(redisModel)
	redis.Engine = redis.Change("master")
	return redis
}

// Change change current engine by name.
// engine load configuration by name from configuration file when cannot found in engine pool.
// given string engine's label name
// returns *redisImpl.Client
func (r *redisModel) Change(name string) *redisImpl.Client {
	var host, port, pass string
	var poolSize int = 1
	envData := NewUtil().GetEnv("CONFIG_PATH")
	resultImpls := NewConfig().SetPath(envData["CONFIG_PATH"]).SetFile("cache").SetPrefix(name+".").GetMulti("host", "port", "pass", "pool_size")
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
		go NewLogger().Error("未找到该标签下的redis配置信息")
		return nil
	}
	engine := redisImpl.NewClient(&redisImpl.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
		PoolSize: poolSize,
	})
	if _, err := engine.Ping(Ctx).Result(); err != nil {
		go NewLogger().Error("(" + name + ")redis连接失败:" + err.Error())
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

// New add new engine by params.
// given string name means label name
// given string host means redis's host
// given string port means redis's port
// given string pass means redis's password
// given int poolSize means conn-pool's size
// returns *redisImpl.Client
func (r *redisModel) New(name, host, port, pass string, poolSize int) *redisImpl.Client {
	engine := redisImpl.NewClient(&redisImpl.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
		PoolSize: poolSize,
	})
	if _, err := engine.Ping(Ctx).Result(); err != nil {
		go NewLogger().Error("(" + name + ")redis连接失败:" + err.Error())
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

// Db change current db in current engine
// given int indev means db's index
// return *redisImpl.Client
func (r *redisModel) Db(index int) *redisImpl.Client {
	r.db = index
	r.Engine = redisImpl.NewClient(&redisImpl.Options{DB: index})
	return r.Engine
}
