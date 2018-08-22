// File file_helper
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

import (
	"os"
)

type FileWriter struct {
	file_name string   // file name
	file      *os.File // file
}

type FileWriterError struct {
	file_name string
	err_info  string
}

func (e *FileWriterError) Error() string { return e.file_name + " error info " + e.err_info }

func NewFileWriter() *FileWriter {
	return &FileWriter{
		file_name: "",
		file:      nil,
	}
}

func (this *FileWriter) Close() {
	if this.file != nil {
		this.file.Sync()
		this.file.Close()
		this.file = nil
	}
}

func (this *FileWriter) Open(name string, truncate bool) bool {
	this.Close()
	this.file_name = name

	open_flag := os.O_CREATE | os.O_RDWR
	if truncate {
		open_flag |= os.O_TRUNC
	} else {
		open_flag |= os.O_APPEND
	}

	file, err := os.OpenFile(name, open_flag, os.ModeExclusive)
	if err != nil {
		return false
	}

	this.file = file

	return true
}

func (this *FileWriter) Reopen(truncate bool) bool {
	return this.Open(this.file_name, truncate)
}

func (this *FileWriter) Write(b []byte) (n int, err error) {
	if this.file != nil {
		return this.file.Write(b)
	}
	return 0, &FileWriterError{this.file_name, "file is nil"}
}

func (this *FileWriter) Flush() {
	if this.file != nil {
		this.file.Sync()
	}
}
