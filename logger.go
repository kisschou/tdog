// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	"os"

	redisImpl "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/**
 * The module for Logger handler.
 * Use Uber's cap and redis message queue.
 * As a replacement of the ELK.
 */
type (
	// logger struct
	logger struct {
	}

	// the writer of writer
	redisWriter struct {
		cli     *redisImpl.Client
		listKey string
	}
)

// newRedisWriter init writer of redis
func newRedisWriter(key string) *redisWriter {
	cli := NewRedis().Engine
	return &redisWriter{
		cli: cli, listKey: key,
	}
}

// Write writes log to the redis
func (w *redisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(Ctx, w.listKey, p).Result()
	return int(n), err
}

// newLogger init zap's Logger module use redis's writer
func newLogger(writer *redisWriter) *zap.Logger {
	// 限制日志输出级别, >= DebugLevel 会打印所有级别的日志
	// 生产环境中一般使用 >= ErrorLevel
	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.DebugLevel
	})

	// use log message of json style
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

// NewLogger init logger module
func NewLogger() *logger {
	return &logger{}
}

// Error writes error message to queue, extends *logger
// given string message of error
func (log *logger) Error(message string) {
	newLogger(newRedisWriter("log:list")).Error(message)
}

// Warn writes warning message to queue, extends *logger
// given string message of warning
func (log *logger) Warn(message string) {
	newLogger(newRedisWriter("log:list")).Warn(message)
}

// Info writes infomation message to queue, extends *logger
// given string message of infomation
func (log *logger) Info(message string) {
	newLogger(newRedisWriter("log:list")).Info(message)
}
