// File file_sink
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

import (
	"fmt"
	"path"
	"strings"
	"time"
)

type FileSinkError struct {
	file_name string
	err_info  string
}

func (e *FileSinkError) Error() string { return e.file_name + " error info " + e.err_info }

// base file sink
type SimpleFileSink struct {
	BaseSink
	file_name   string      // file name
	file_writer *FileWriter // file writer
}

func NewSimpleFileSink(filename string, thread_safe bool) *SimpleFileSink {
	file_sink := &SimpleFileSink{
		file_name:   filename,
		file_writer: NewFileWriter(),
	}

	file_sink.DefaultInit()

	if thread_safe {
		file_sink.SetLocker(NewExclusiveLocker())
	} else {
		file_sink.SetLocker(NewNullLocker())
	}

	file_sink.file_writer.Open(filename, false)

	return file_sink
}

// create a simple file sink of thread no safe(in single thread situation)
func NewSimpleFileSinkSt(filename string) *SimpleFileSink {
	return NewSimpleFileSink(filename, false)
}

// create a simple file sink of thread safe
func NewSimpleFileSinkMt(filename string) *SimpleFileSink {
	return NewSimpleFileSink(filename, true)
}

// override the Log method
func (this *SimpleFileSink) Log(msg string) {
	this.file_writer.Write([]byte(msg))
}

// override the Flush method
func (this *SimpleFileSink) Flush() {
	this.file_writer.Flush()
}

// rorating file sink
type RotatingFileSink struct {
	BaseSink
	base_file_name      string      // base file name
	max_file_size       int64       // max file size each log file
	max_files           int32       // max file count
	cur_file_size       int64       // current file size
	file_writer         *FileWriter // file writer
	base_file_name_base string
	base_file_name_ext  string
}

func NewRotatingFileSink(
	_base_file_name string, _max_size int64, _max_files int32, thread_safe bool) *RotatingFileSink {
	sink := &RotatingFileSink{
		base_file_name: _base_file_name,
		max_file_size:  _max_size,
		max_files:      _max_files,
		file_writer:    NewFileWriter(),
		cur_file_size:  0,
	}

	sink.DefaultInit()

	if thread_safe {
		sink.SetLocker(NewExclusiveLocker())
	} else {
		sink.SetLocker(NewNullLocker())
	}

	sink.base_file_name_ext = path.Ext(sink.base_file_name)
	sink.base_file_name_base = strings.TrimSuffix(sink.base_file_name, sink.base_file_name_ext)

	first_file_name := sink.calc_file_name(0)
	sink.file_writer.Open(first_file_name, false)

	return sink
}

// create not thread safe rotating logger
func NewRotatingFileSinkSt(
	_base_file_name string, _max_size int64, _max_files int32) *RotatingFileSink {
	return NewRotatingFileSink(_base_file_name, _max_size, _max_files, false)
}

// create thread safe rotating logger
func NewRotatingFileSinkMt(
	_base_file_name string, _max_size int64, _max_files int32) *RotatingFileSink {
	return NewRotatingFileSink(_base_file_name, _max_size, _max_files, true)
}

// calc rorating file name
func (this *RotatingFileSink) calc_file_name(file_index int32) string {
	if file_index != 0 {
		return fmt.Sprintf("%s.%d%s", this.base_file_name_base, file_index, this.base_file_name_ext)
	} else {
		return this.base_file_name
	}
}

// rorating the files
// Rotate files:
// log.txt -> log.1.txt
// log.1.txt -> log.2.txt
// log.2.txt -> log.3.txt
// log.3.txt -> delete
func (this *RotatingFileSink) rotate() {
	this.file_writer.Close()

	for i := this.max_files; i > 0; i-- {
		src_file_name := this.calc_file_name(i - 1)
		tar_file_name := this.calc_file_name(i)

		tar_file_exit, _ := PathExists(tar_file_name)
		if tar_file_exit {
			// remove target file
			RemoveFile(tar_file_name)
		}

		src_file_exit, _ := PathExists(src_file_name)
		if src_file_exit {
			RenamePath(src_file_name, tar_file_name)
		}
	}

	// reopen the file writer
	this.file_writer.Reopen(true)
}

// override the Log method
func (this *RotatingFileSink) Log(msg string) {
	msg_bytes := []byte(msg)

	this.cur_file_size += (int64)(len(msg_bytes))
	if this.cur_file_size > this.max_file_size {
		this.rotate()
		this.cur_file_size = (int64)(len(msg_bytes))
	}
	this.file_writer.Write(msg_bytes)
}

// override the Flush method
func (this *RotatingFileSink) Flush() {
	this.file_writer.Flush()
}

// daily file sink
type DailyFileSink struct {
	BaseSink
	base_file_name      string      // base file name
	rotation_hour       int32       // roration hour
	rotation_minute     int32       // rotation minute
	file_writer         *FileWriter // file writer
	base_file_name_base string
	base_file_name_ext  string
	next_rotation_time  time.Time
}

// create new DailyFileSink
func NewDailyFileSink(_base_file_name string, _rotation_hour int32, _rotation_minute int32, thread_safe bool) (*DailyFileSink, error) {

	if _rotation_hour < 0 || _rotation_hour > 23 || _rotation_minute < 0 || _rotation_minute > 59 {
		error_info := fmt.Sprintf("rotation hour[%d] or rotaion minute[%d] invalid", _rotation_hour, _rotation_minute)
		return nil, &FileSinkError{_base_file_name, error_info}
	}

	sink := &DailyFileSink{
		base_file_name:  _base_file_name,
		rotation_hour:   _rotation_hour,
		rotation_minute: _rotation_minute,
		file_writer:     NewFileWriter(),
	}

	sink.DefaultInit()

	if thread_safe {
		sink.SetLocker(NewExclusiveLocker())
	} else {
		sink.SetLocker(NewNullLocker())
	}

	// calc next rotation time
	sink.next_rotation_time = sink.calc_next_rotation_time()

	sink.base_file_name_ext = path.Ext(sink.base_file_name)
	sink.base_file_name_base = strings.TrimSuffix(sink.base_file_name, sink.base_file_name_ext)

	// calc current file name
	sink.file_writer.Open(sink.calc_file_name(), false)

	return sink, nil
}

// create not thread safe daily file sink
func NewDailyFileSinkSt(_base_file_name string, _rotation_hour int32, _rotation_minute int32) (*DailyFileSink, error) {
	return NewDailyFileSink(_base_file_name, _rotation_hour, _rotation_minute, false)
}

// create thread safe daily file sink
func NewDailyFileSinkMt(_base_file_name string, _rotation_hour int32, _rotation_minute int32) (*DailyFileSink, error) {
	return NewDailyFileSink(_base_file_name, _rotation_hour, _rotation_minute, true)
}

// calc file name
func (this *DailyFileSink) calc_file_name() string {
	now := time.Now()
	return fmt.Sprintf("%s_%04d-%02d-%02d%s",
		this.base_file_name_base, now.Year(), now.Month(), now.Day(), this.base_file_name_ext)
}

// calc next rotation time
func (this *DailyFileSink) calc_next_rotation_time() time.Time {
	now := time.Now()
	next_rotation := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		int(this.rotation_hour),
		int(this.rotation_minute),
		int(0),
		int(0),
		now.Location())

	if next_rotation.Before(now) {
		next_rotation = next_rotation.AddDate(0, 0, 1)
	}

	return next_rotation
}

// override the Log method
func (this *DailyFileSink) Log(msg string) {

	// check next create a new file
	if time.Now().After(this.next_rotation_time) {
		this.file_writer.Open(this.calc_file_name(), false)
		this.next_rotation_time = this.calc_next_rotation_time()
	}
	this.file_writer.Write([]byte(msg))
}

// override the Flush method
func (this *DailyFileSink) Flush() {
	this.file_writer.Flush()
}

// hourly file sink
type HourlyFileSink struct {
	BaseSink
	base_file_name      string      // base file name
	file_writer         *FileWriter // file writer
	base_file_name_base string
	base_file_name_ext  string
	next_rotation_time  time.Time
}

// new hourly file sink
func NewHourlyFileSink(_base_file_name string, thread_safe bool) *HourlyFileSink {
	sink := &HourlyFileSink{
		base_file_name: _base_file_name,
		file_writer:    NewFileWriter(),
	}

	sink.DefaultInit()

	if thread_safe {
		sink.SetLocker(NewExclusiveLocker())
	} else {
		sink.SetLocker(NewNullLocker())
	}

	// calc next rotation time
	sink.next_rotation_time = sink.calc_next_rotation_time()

	sink.base_file_name_ext = path.Ext(sink.base_file_name)
	sink.base_file_name_base = strings.TrimSuffix(sink.base_file_name, sink.base_file_name_ext)

	// calc current file name
	sink.file_writer.Open(sink.calc_file_name(), false)

	return sink
}

// create not thread safe hourly file sink
func NewHourlyFileSinkSt(_base_file_name string) *HourlyFileSink {
	return NewHourlyFileSink(_base_file_name, false)
}

// create thread safe hourly file sink
func NewHourlyFileSinkMt(_base_file_name string) *HourlyFileSink {
	return NewHourlyFileSink(_base_file_name, true)
}

// calc next rotation time
func (this *HourlyFileSink) calc_next_rotation_time() time.Time {
	now := time.Now()
	next_rotation := now.Add(time.Hour)

	return next_rotation
}

// calc next file name
func (this *HourlyFileSink) calc_file_name() string {
	now := time.Now()
	return fmt.Sprintf("%s_%04d-%02d-%02d-%02d%s",
		this.base_file_name_base, now.Year(), now.Month(), now.Day(), now.Hour(), this.base_file_name_ext)
}

// override the Log method
func (this *HourlyFileSink) Log(msg string) {

	// check next create a new file
	if time.Now().After(this.next_rotation_time) {
		this.file_writer.Open(this.calc_file_name(), false)
		this.next_rotation_time = this.calc_next_rotation_time()
	}
	this.file_writer.Write([]byte(msg))
}

// override the Flush method
func (this *HourlyFileSink) Flush() {
	this.file_writer.Flush()
}
