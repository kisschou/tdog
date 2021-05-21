// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

/**
 * The module for MySQL handler.
 *
 * @Author: Kisschou
 * @Build: 2021-04-22
 */
type (
	mySql struct {
		engineList map[string]*xorm.Engine // engine pool
		Engine     *xorm.Engine            // current engine
	}

	MySqlConf struct {
		engine       string
		Host         string
		Port         string
		User         string
		Pass         string
		Db           string
		Charset      string
		Prefix       string
		dsn          string
		Debug        bool // is start debug
		MaxIdleConns int  // 连接池的空闲数大小
		MaxOpenConns int  // max connections count
	}
)

// NewMySQL init MySQL module
// auto-load master configuration and make it actived.
// so...must write master configuration of mysql into configuration.
func NewMySQL() *mySql {
	sql := new(mySql)
	sql.Engine = sql.Change("master")
	return sql
}

// Change change current engine by name.
// engine load configuration by name from configuration file when cannot found in engine pool.
// given string engine's label name
// returns *xorm.Engine
func (sql *mySql) Change(name string) *xorm.Engine {
	engine := sql.engineList[name]
	if engine == nil {
		var err error
		conf := loadConf(name)
		engine, err = xorm.NewEngine(conf.engine, conf.dsn)
		if err != nil {
			go NewLogger().Error(err.Error())
			return nil
		}

		// 日志打印SQL
		engine.ShowSQL(conf.Debug) // 设置连接池的空闲数大小 engine.SetMaxIdleConns(conf.maxIdleConns)
		// 设置最大连接数
		engine.SetMaxOpenConns(conf.MaxOpenConns)
		// 名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
		engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, conf.Prefix))
		// 加入连接池
		if sql.engineList == nil {
			sql.engineList = make(map[string]*xorm.Engine, 0)
		}
		sql.engineList[name] = engine
	}
	return engine
}

// New add new engine by params.
// given string name means label name
// given *MySqlConf conf means mysql's config
// returns *xorm.Engine
func (sql *mySql) New(name string, conf *MySqlConf) *xorm.Engine {
	dsn := conf.User + ":" + conf.Pass + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Db + "?charset=" + conf.Charset + "&parseTime=True&loc=Local"
	debug := NewConfig().Get("database.debug").ToBool()
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		go NewLogger().Error(err.Error())
		return nil
	}

	// 日志打印SQL
	engine.ShowSQL(debug) // 设置连接池的空闲数大小 engine.SetMaxIdleConns(conf.maxIdleConns)
	// 设置最大连接数
	engine.SetMaxOpenConns(5)
	// 名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
	engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, conf.Prefix))
	// 加入连接池
	sql.engineList[name] = engine

	return sql.Change(name)
}

// loadConf load configuration from configuration file.
// write to new *mysqlConf and return it.
func loadConf(name string) *MySqlConf {
	envData := NewUtil().GetEnv("CONFIG_PATH")
	configResults := NewConfig().SetPath(envData["CONFIG_PATH"]).SetFile("database").SetPrefix(name+".").GetMulti("host", "port", "user", "pass", "db", "charset", "prefix")
	return &MySqlConf{
		engine:       "mysql",
		Host:         configResults["host"].ToString(),
		Port:         configResults["port"].ToString(),
		User:         configResults["user"].ToString(),
		Pass:         configResults["pass"].ToString(),
		Db:           configResults["db"].ToString(),
		Charset:      configResults["charset"].ToString(),
		Prefix:       configResults["prefix"].ToString(),
		dsn:          configResults["user"].ToString() + ":" + configResults["pass"].ToString() + "@tcp(" + configResults["host"].ToString() + ":" + configResults["port"].ToString() + ")/" + configResults["db"].ToString() + "?charset=" + configResults["charset"].ToString() + "&parseTime=True&loc=Local",
		Debug:        NewConfig().Get("database.debug").ToBool(),
		MaxIdleConns: 5,
		MaxOpenConns: 5,
	}
}
