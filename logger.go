package tdog

import (
	"os"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger struct {
	}

	// 为 logger 提供写入 redis 队列的 io 接口
	redisWriter struct {
		cli     *redis.Client
		listKey string
	}
)

func (log *Logger) NewRedisWriter(key string) *redisWriter {
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	return &redisWriter{
		cli: cli, listKey: key,
	}
}

func (w *redisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(w.listKey, p).Result()
	return int(n), err
}

func (log *Logger) NewLogger(writer *redisWriter) *zap.Logger {
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

func (log *Logger) Error(message string) {
	log.NewLogger(log.NewRedisWriter("log:list")).Error(message)
}

func (log *Logger) Warn(message string) {
	log.NewLogger(log.NewRedisWriter("log:list")).Warn(message)
}

func (log *Logger) Info(message string) {
	log.NewLogger(log.NewRedisWriter("log:list")).Info(message)
}
