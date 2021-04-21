package tdog

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

type (
	MySql struct {
		IsCluster   bool
		Engine      *xorm.Engine
		EngineGroup *xorm.EngineGroup
	}

	MySqlConf struct {
		engine  string
		host    string
		port    string
		user    string
		pass    string
		db      string
		charset string
		prefix  string
		dsn     string
		debug   bool
	}
)

func (mysql *MySql) NewEngine() {
	var err error

	// master-slave mode.
	mysql.IsCluster = NewConfig().Get("database.master_slave").ToBool()
	if mysql.IsCluster {
		mysql.EngineGroup, err = cluster()
	} else {
		mysql.Engine, err = singleton()
	}

	if err != nil {
		NewLogger().Error(err.Error())
		os.Exit(0)
	}
}

func cluster() (engineGroup *xorm.EngineGroup, err error) {
	conns := []string{}
	// append master conf.
	var myCnf *MySqlConf
	myCnf = myCnf.Get("master")
	conns = append(conns, myCnf.dsn)

	// append slave conf.
	slaveList := NewConfig().Get("database.slave_list").ToStringSlice()
	for _, slave := range slaveList {
		myCnf = myCnf.Get(slave)
		conns = append(conns, myCnf.dsn)
	}

	/**
	 * 负载策略
	 *
	 * 1. 随机访问负载策略 xorm.RandomPolicy()
	 * 2. 权重随机访问负载策略 xorm.WeightRandomPolicy([]int{2, 3})
	 * 3. 轮询访问负载策略 xorm.RoundRobinPolicy()
	 * 4. 权重轮询访问负载策略 xorm.WeightRoundRobinPolicy([]int{2, 3})
	 * 5. 最小连接数访问负载策略 xorm.LeastConnPolicy()
	 * 6. 自定义
	 *		type GroupPolicy interface {
	 *			Slave(*EngineGroup) *Engine
	 *		}
	 *
	 * 注: 当前采用随机访问负载策略
	 */
	engineGroup, err = xorm.NewEngineGroup(myCnf.engine, conns, xorm.RandomPolicy())

	// 日志打印SQL
	engineGroup.ShowSQL(myCnf.debug)
	// 设置连接池的空闲数大小
	engineGroup.SetMaxIdleConns(5)
	// 设置最大连接数
	engineGroup.SetMaxOpenConns(5)
	// 名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
	engineGroup.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, myCnf.prefix))

	return
}

func singleton() (engine *xorm.Engine, err error) {
	var myCnf *MySqlConf
	myCnf = myCnf.Get("master")
	engine, err = xorm.NewEngine(myCnf.engine, myCnf.dsn)

	if err != nil {
		return
	}

	// 日志打印SQL
	engine.ShowSQL(myCnf.debug)
	// 设置连接池的空闲数大小
	engine.SetMaxIdleConns(5)
	// 设置最大连接数
	engine.SetMaxOpenConns(5)
	// 名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
	engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, myCnf.prefix))

	return
}

func (mysql *MySql) Ping() bool {
	if err := mysql.Engine.Ping(); err != nil {
		NewLogger().Error(err.Error())
		return false
	}
	return true
}

func (myCnf *MySqlConf) Get(hostType string) *MySqlConf {
	mysqlUser := NewConfig().Get("database.master.user").ToString()
	mysqlPass := NewConfig().Get("database." + hostType + ".pass").ToString()
	mysqlHost := NewConfig().Get("database." + hostType + ".host").ToString()
	mysqlPort := NewConfig().Get("database." + hostType + ".port").ToString()
	mysqlDb := NewConfig().Get("database." + hostType + ".db").ToString()
	mysqlCharset := NewConfig().Get("database." + hostType + ".charset").ToString()
	mysqlPrefix := ""
	impl := NewConfig().Get("database." + hostType + ".prefix")
	if impl.IsExists() {
		mysqlPrefix = impl.ToString()
	}
	dsn := mysqlUser + ":" + mysqlPass + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDb + "?charset=" + mysqlCharset + "&parseTime=True&loc=Local"
	debug := NewConfig().Get("database.debug").ToBool()

	return &MySqlConf{
		engine:  "mysql",
		host:    mysqlHost,
		port:    mysqlPort,
		user:    mysqlUser,
		pass:    mysqlPass,
		db:      mysqlDb,
		charset: mysqlCharset,
		prefix:  mysqlPrefix,
		dsn:     dsn,
		debug:   debug,
	}
}
