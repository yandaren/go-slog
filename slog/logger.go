// File logger.go
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

import (
	"bytes"
	"container/list"
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	log_lvl           LogLevel   // the log level
	print_millseconds bool       // if print milliseconds
	print_lvl         bool       // if print lvl
	msg_content_blank bool       // if log content and log head has blank
	log_name          string     // the logger name
	formatter_pattern string     // log formatter pattern
	sink_list         *list.List // sink list
	line_sperator     string     // line seprator
}

func NewLogger(name string) *Logger {
	logger := &Logger{
		log_lvl:           LvlDebug,
		print_millseconds: true,
		print_lvl:         true,
		msg_content_blank: true,
		log_name:          name,
		sink_list:         list.New(),
	}

	if strings.ToLower(runtime.GOOS) == "windows" {
		logger.line_sperator = "\r\n"
	} else {
		logger.line_sperator = "\n"
	}

	return logger
}

func (this *Logger) SetPrintLvl(p bool) {
	this.print_lvl = p
}

func (this *Logger) IsPrintLvl() bool {
	return this.print_lvl
}

func (this *Logger) SetPrintMilliseconds(p bool) {
	this.print_millseconds = p
}

func (this *Logger) IsPrintMilliseconds() bool {
	return this.print_millseconds
}

func (this *Logger) SetMsgContenHasBlank(p bool) {
	this.msg_content_blank = p
}

func (this *Logger) IsMsgContentHasBlank() bool {
	return this.msg_content_blank
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

	var buffer bytes.Buffer
	// time
	buffer.WriteString(fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()))

	if this.IsPrintMilliseconds() {
		buffer.WriteString(fmt.Sprintf(".%03d]", now.Nanosecond()/1000000))
	} else {
		buffer.WriteString("]")
	}

	if this.IsPrintLvl() {
		buffer.WriteString(fmt.Sprintf("[%-5s]", lvl.String()))
	}

	if this.IsMsgContentHasBlank() {
		buffer.WriteString(" ")
	}

	buffer.WriteString(msg_content)
	buffer.WriteString(this.line_sperator)

	msg := buffer.String()

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

func (this *Logger) Flush() {
	for e := this.sink_list.Front(); e != nil; e = e.Next() {
		sink := e.Value.(Sink)
		if sink != nil {

			// lock
			sink.Lock()

			sink.Flush()

			// unlock
			sink.Unlock()
		}
	}
}
