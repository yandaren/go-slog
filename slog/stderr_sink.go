// File stderr_sink.go
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

import (
	"fmt"
	"os"
)

type StderrSink struct {
	BaseSink
}

func NewStderrSink(thread_safe bool) *StderrSink {
	sink := &StderrSink{}
	sink.DefaultInit()

	if thread_safe {
		sink.SetLocker(NewExclusiveLocker())
	} else {
		sink.SetLocker(NewNullLocker())
	}

	return sink
}

// create not thread safe stderr sink
func NewStderrSinkSt() *StderrSink {
	return NewStderrSink(false)
}

// create thread safe stderr sink
func NewStderrSinkMt() *StderrSink {
	return NewStderrSink(true)
}

func (this *StderrSink) Log(msg string) {
	fmt.Fprintf(os.Stderr, "%s", msg)
}

func (this *StderrSink) Flush() {
}
