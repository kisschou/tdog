package tdog

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	Path  string
	File  string
	Level string
}

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

func (log *Logger) Start() {
	UtilTdog := new(Util)
	if !UtilTdog.checkPortAlived(7999) {
		mux := http.NewServeMux()
		mux.HandleFunc("/error", log.Error)
		mux.HandleFunc("/warning", log.Warning)
		mux.HandleFunc("/access", log.Access)
		go http.ListenAndServe(":7999", mux)
	}
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

func (log *Logger) Writer(context string) {
	fileName := log.Path + log.File
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("write to log file FAIL: ", err)
		os.Exit(0)
	}

	// 实例化
	LogImpl := logrus.New()
	// 设置输出
	LogImpl.Out = src
	// 设置输出日志中添加文件名和方法信息
	LogImpl.SetReportCaller(true)
	// 设置日志格式
	LogImpl.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 时间格式化
	})
	/**
	 * 日志的级别（来自@dylanbeattie）
	 * - Fatal：网站挂了，或者极度不正常
	 * - Error：跟遇到的用户说对不起，可能有bug
	 * - Warn：记录一下，某事又发生了
	 * - Info：提示一切正常
	 * - debug：没问题，就看看堆栈
	 */
	switch log.Level {
	case "error":
		LogImpl.SetLevel(logrus.ErrorLevel)
		LogImpl.Error("%s", context)
		break
	case "warning":
		LogImpl.SetLevel(logrus.WarnLevel)
		LogImpl.Warn("%s", context)
		break
	case "access":
		LogImpl.SetLevel(logrus.InfoLevel)
		LogImpl.Info("%s", context)
		break
	default:
		LogImpl.SetLevel(logrus.InfoLevel)
		LogImpl.Infof("%s", context)
		break
	}
}
