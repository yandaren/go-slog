// File stdout_sink.go
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

import (
	"fmt"
)

type StdoutSink struct {
	BaseSink
}

func NewStdoutSink(thread_safe bool) *StdoutSink {
	sink := &StdoutSink{}
	sink.DefaultInit()

	if thread_safe {
		sink.SetLocker(NewExclusiveLocker())
	} else {
		sink.SetLocker(NewNullLocker())
	}

	return sink
}

// create not thread safe stdout sink
func NewStdoutSinkSt() *StdoutSink {
	return NewStdoutSink(false)
}

// create thread safe stdout sink
func NewStdoutSinkMt() *StdoutSink {
	return NewStdoutSink(true)
}

func (this *StdoutSink) Log(msg string) {
	fmt.Printf("%s", msg)
}

func (this *StdoutSink) Flush() {
}
