// File logger.go
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

import (
	"container/list"
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	log_lvl           LogLevel   // the log level
	log_name          string     // the logger name
	formatter_pattern string     // log formatter pattern
	sink_list         *list.List // sink list
	line_sperator     string     // line seprator
}

func NewLogger(name string) *Logger {
	logger := &Logger{
		log_lvl:   LvlDebug,
		log_name:  name,
		sink_list: list.New(),
	}

	if strings.ToLower(runtime.GOOS) == "windows" {
		logger.line_sperator = "\r\n"
	} else {
		logger.line_sperator = "\n"
	}

	return logger
}

func (this *Logger) Name() string {
	return this.log_name
}

func (this *Logger) AppendSink(sink Sink) *Logger {
	this.sink_list.PushBack(sink)
	return this
}

func (this *Logger) SetLogLvl(lvl LogLevel) {
	this.log_lvl = lvl
}

func (this *Logger) GetLogLvl() LogLevel {
	return this.log_lvl
}

func (this *Logger) ShouldLog(lvl LogLevel) bool {
	return lvl >= this.log_lvl
}

func (this *Logger) log_msg(lvl LogLevel, format string, args ...interface{}) {
	if !this.ShouldLog(lvl) {
		return
	}

	msg_content := fmt.Sprintf(format, args...)

	// time prefix
	now := time.Now()
	msg := fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d.%03d][%-5s] %s%s",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000, lvl.String(), msg_content, this.line_sperator)

	for e := this.sink_list.Front(); e != nil; e = e.Next() {
		sink := e.Value.(Sink)
		if sink != nil && sink.ShoudLog(lvl) {

			// lock
			sink.Lock()

			sink.Log(msg)
			if sink.IsForceFlush() {
				sink.Flush()
			}

			// unlock
			sink.Unlock()
		}
	}
}

func (this *Logger) Debug(format string, args ...interface{}) {
	this.log_msg(LvlDebug, format, args...)
}

func (this *Logger) Info(format string, args ...interface{}) {
	this.log_msg(LvlInfo, format, args...)
}

func (this *Logger) Warn(format string, args ...interface{}) {
	this.log_msg(LvlWarn, format, args...)
}

func (this *Logger) Error(format string, args ...interface{}) {
	this.log_msg(LvlError, format, args...)
}

func (this *Logger) Fatal(format string, args ...interface{}) {
	this.log_msg(LvlFatal, format, args...)

}
