package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger struct {
		Path  string
		File  string
		Level string
	}

	// 为 logger 提供写入 redis 队列的 io 接口
	redisWriter struct {
		cli     *redis.Client
		listKey string
	}
)

func (log *Logger) BuildFilePath() *Logger {
	var util *Util
	filePath, _ := os.Getwd()
	filePath += "/runtime/" + time.Now().Format("2006-01-02") + "/"
	util.DirExistsAndCreate(filePath)
	log.Path = filePath
	return log
}

func (log *Logger) BuildFileName(fileName string) *Logger {
	log.Level = fileName
	log.File = fileName + ".log"
	return log
}

func (log *Logger) Error(res http.ResponseWriter, req *http.Request) {
	log.BuildFilePath().BuildFileName("error")
}

func (log *Logger) Warning(res http.ResponseWriter, req *http.Request) {
	log.BuildFilePath().BuildFileName("warning")
}

func (log *Logger) Access(res http.ResponseWriter, req *http.Request) {
	log.BuildFilePath().BuildFileName("access").Writer("")
}

func NewRedisWriter(key string, cli *redis.Client) *redisWriter {
	return &redisWriter{
		cli: cli, listKey: key,
	}
}

func (w *redisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(w.listKey, p).Result()
	return int(n), err
}

func NewLogger(writer *redisWriter) *zap.Logger {
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

func main() {
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	writer := NewRedisWriter("log_list", cli)
	logger := NewLogger(writer)

	logger.Info("test logger info", zap.String("hello", "logger"))
	logger.Error("test logger info", zap.String("hello", "logger"))
	logger.Warn("test logger info", zap.String("hello", "logger"))
}
