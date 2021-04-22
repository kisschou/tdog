package tdog

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

type (
	mySql struct {
		engineList    map[string]*xorm.Engine // 引擎列表
		defaultEngine string                  // 默认引擎名
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

func NewMySQL() *mySql {
	return &mySql{}
}

func (sql *mySql) Engine(name string) *xorm.Engine {
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
		engine.ShowSQL(conf.debug)
		// 设置连接池的空闲数大小
		engine.SetMaxIdleConns(conf.maxIdleConns)
		// 设置最大连接数
		engine.SetMaxOpenConns(conf.maxOpenConns)
		// 名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
		engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, conf.prefix))
		// 加入连接池
		sql.engineList[name] = engine
	}
	return engine
}

func (sql *mySql) NewConn() {
}

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
