package tdog

import (
	"os"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/**
 * 日志基础模块
 * 采用Uber开源的zap开发，采用redis消息队列, 替代业内ELK.
 */
type (
	Logger struct {
	}

	// 为 logger 提供写入 redis 队列的 io 接口
	redisWriter struct {
		cli     *redis.Client
		listKey string
	}
)

func newRedisWriter(key string) *redisWriter {
	cli := redis.NewClient(&redis.Options{
		Addr: NewConfig().Get("cache.master.host").ToString() + ":" + NewConfig().Get("cache.master.port").ToString(),
	})
	return &redisWriter{
		cli: cli, listKey: key,
	}
}

func (w *redisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(w.listKey, p).Result()
	return int(n), err
}

func newLogger(writer *redisWriter) *zap.Logger {
	// 限制日志输出级别, >= DebugLevel 会打印所有级别的日志
	// 生产环境中一般使用 >= ErrorLevel
	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.DebugLevel
	})

	// 使用 JSON 格式日志
	jsonEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	stdCore := zapcore.NewCore(jsonEnc, zapcore.Lock(os.Stdout), lowPriority)

	// addSync 将 io.Writer 装饰为 WriteSyncer
	// 故只需要一个实现 io.Writer 接口的对象即可
	syncer := zapcore.AddSync(writer)
	redisCore := zapcore.NewCore(jsonEnc, syncer, lowPriority)

	// 集成多个 core
	core := zapcore.NewTee(stdCore, redisCore)

	// logger 输出到 console 且标识调用代码行
	return zap.New(core).WithOptions(zap.AddCaller())
}

/**
 * 初始化日志模块
 *
 * @return *Logger
 */
func NewLogger() *Logger {
	return &Logger{}
}

/**
 * 输出错误日志
 *
 * @param string message 消息内容
 *
 * @return nil
 */
func (log *Logger) Error(message string) {
	newLogger(newRedisWriter("log:list")).Error(message)
}

/**
 * 输出警告日志
 *
 * @param string message 消息内容
 *
 * @return nil
 */
func (log *Logger) Warn(message string) {
	newLogger(newRedisWriter("log:list")).Warn(message)
}

/**
 * 输出消息日志
 *
 * @param string message 消息内容
 *
 * @return nil
 */
func (log *Logger) Info(message string) {
	newLogger(newRedisWriter("log:list")).Info(message)
}
