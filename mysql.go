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

	mysqlConf struct {
		engine       string
		host         string
		port         string
		user         string
		pass         string
		db           string
		charset      string
		prefix       string
		dsn          string
		debug        bool // is start debug
		maxIdleConns int  // 连接池的空闲数大小
		maxOpenConns int  // max connections count
	}
)

// NewMySQL init MySQL module
// auto-load master configuration and make it actived
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
		engine.ShowSQL(conf.debug) // 设置连接池的空闲数大小 engine.SetMaxIdleConns(conf.maxIdleConns)
		// 设置最大连接数
		engine.SetMaxOpenConns(conf.maxOpenConns)
		// 名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
		engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, conf.prefix))
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
// given string host means mysql's host
// given string port means mysql's port
// given string user means mysql's user
// given string pass means mysql's password
// given string db means mysql's db
// given string db means mysql's charset
// given string prefix means mysql's prefix
// returns *xorm.Engine
func (sql *mySql) New(name, host, port, user, pass, db, charset, prefix string) *xorm.Engine {
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + db + "?charset=" + charset + "&parseTime=True&loc=Local"
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
	engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, prefix))
	// 加入连接池
	sql.engineList[name] = engine

	return sql.Change(name)
}

// loadConf load configuration from configuration file.
// write to new *mysqlConf and return it.
func loadConf(name string) *mysqlConf {
	configResults := NewConfig().SetFile("database").SetPrefix(name+".").GetMulti("host", "port", "user", "pass", "db", "charset", "prefix")
	return &mysqlConf{
		engine:       "mysql",
		host:         configResults["host"].ToString(),
		port:         configResults["port"].ToString(),
		user:         configResults["user"].ToString(),
		pass:         configResults["pass"].ToString(),
		db:           configResults["db"].ToString(),
		charset:      configResults["charset"].ToString(),
		prefix:       configResults["prefix"].ToString(),
		dsn:          configResults["user"].ToString() + ":" + configResults["pass"].ToString() + "@tcp(" + configResults["host"].ToString() + ":" + configResults["port"].ToString() + ")/" + configResults["db"].ToString() + "?charset=" + configResults["charset"].ToString() + "&parseTime=True&loc=Local",
		debug:        NewConfig().Get("database.debug").ToBool(),
		maxIdleConns: 5,
		maxOpenConns: 5,
	}
}
