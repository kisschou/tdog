package tdog

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

/**
 * MySQL模块
 *
 * @Author: Kisschou
 * @Build: 2021-04-22
 */
type (
	mySql struct {
		engineList map[string]*xorm.Engine // 引擎列表
		Engine     *xorm.Engine            // 默认引擎
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
		debug        bool // 是否开启调试
		maxIdleConns int  // 连接池的空闲数大小
		maxOpenConns int  // 最大连接数
	}
)

/**
 * 初始化MySQL模块
 *
 * @return *xorm.Engine
 */
func NewMySQL() *mySql {
	sql := new(mySql)
	sql.Engine = sql.Change("master")
	return sql
}

/**
 * 引擎切换
 *
 * @param string name 切换的引擎名
 *
 * @return *xorm.Engine
 */
func (sql *mySql) Change(name string) *xorm.Engine {
	engine := sql.engineList[name]
	if engine == nil {
		var err error
		conf := loadConf(name)
		engine, err = xorm.NewEngine(conf.engine, conf.dsn)
		if err != nil {
			NewLogger().Error(err.Error())
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

/**
 * 新增其他mysql引擎
 *
 * @param string name 引擎名
 * @param string host 连接地址
 * @param string user 账户
 * @param string pass 密码
 * @param string db 数据库名
 * @param string charset 字符集
 * @param string prefix 表前缀
 *
 * @return *xorm.Engine
 */
func (sql *mySql) New(name, host, port, user, pass, db, charset, prefix string) *xorm.Engine {
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + db + "?charset=" + charset + "&parseTime=True&loc=Local"
	debug := NewConfig().Get("database.debug").ToBool()
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		NewLogger().Error(err.Error())
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

/**
 * 从配置文件中加载指定的mysql配置
 *
 * @param string name 数据库配置类名
 *
 * @return *mysqlConf
 */
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
