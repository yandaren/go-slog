package slog_example

import (
	"fmt"
	"go-slog/slog"
	"time"
)

func Slog_test() {
	for _, lg_lvl := range slog.AllLogLevels {
		fmt.Printf("%d = %s\n", lg_lvl, lg_lvl.String())
	}

	fmt.Printf("----------------------------------\n")
	var lvl_strings = []string{
		"debug",
		"info",
		"warn",
		"warning",
		"error",
		"fatal",
		"none",
		"sblv",
	}

	for _, lg_lvl_str := range lvl_strings {
		lg_lvl, err := slog.ParseLevel(lg_lvl_str)
		if err != nil {
			fmt.Printf("parse lg_lvl_str[%s] error[%s]\n", lg_lvl_str, err.Error())
		} else {
			fmt.Printf("log_lvl_str[%s] = %d\n", lg_lvl_str, lg_lvl)
		}
	}

	fmt.Printf("---------------slog test---------------\n")
	logger := slog.NewStdoutLoggerSt("stdout_logger")

	logger.Debug("slog stdoutlogger test")
	logger.Info("slog stdoutlogger test")
	logger.Warn("slog stdoutlogger test, %d", 3)

	stderr_logger := slog.NewStderrLoggerSt("stderr_logger")
	stderr_logger.Debug("slog stderr_logger test")
	stderr_logger.Info("slog stderr_logger test")
	stderr_logger.Warn("slog stderr_logger test, %d", 3)
}

func File_writer_test() {
	fmt.Printf("file_writer_test\n")

	fwriter := slog.NewFileWriter()

	file_name := "file_writer_test.txt"

	if !fwriter.Open(file_name, false) {
		fmt.Printf("create file[%s] failed\n", file_name)
		return
	}

	defer fwriter.Close()

	for i := 0; i < 10; i++ {
		fwriter.Write([]byte("11111111111111111111\n"))
	}

	file_name1 := "file_writer_test1.txt"

	if !fwriter.Open(file_name1, false) {
		fmt.Printf("create file[%s] failed\n", file_name1)
		return
	}

	defer fwriter.Close()

	for i := 0; i < 10; i++ {
		fwriter.Write([]byte("2222222222222222222222\n"))
	}
}

// 这个就是最简单的单个文件的logger
func Simple_file_logger_test() {
	fmt.Printf("simple_file_logger_test\n")

	logger := slog.NewBasicLoggerSt("base_logger", "basic_logger.txt")

	for i := 0; i < 10; i++ {
		logger.Debug("base logger debug log test")
		logger.Info("base logger info log test")
		logger.Warn("base logger warn log test")
	}
}

// 这个指定每个文件的最大大小，以及维护的最大文件格式
// 当前文件大小到达指定的最大文件大小之后，就会进行一次rotating
// 比如最多保留3个文件的话
// Rotate files:
// log.txt -> log.1.txt
// log.1.txt -> log.2.txt
// log.2.txt -> log.3.txt
// log.3.txt -> delete
func Rotating_logger_test() {
	fmt.Printf("rotating_logger_test\n")

	logger := slog.NewRotatingLoggerSt("rotating_logger", "rotating_logger.txt", 500, 5)

	for i := 0; i < 20; i++ {
		logger.Debug("rorating msg xxx now_time[%s]", time.Now().String())
	}
}

// 这个是每日一个文件的logger
func Daily_logger_test() {
	fmt.Printf("daily_logger_test\n")

	logger := slog.NewDailyLoggerSt("daily_logger", "daily_logger.txt", 12, 30)

	for i := 0; i < 20; i++ {
		logger.Debug("daily_logger test")
	}
}

// 这个是每个小时会重新生成一个文件的logger
func Hourly_logger_test() {
	fmt.Printf("Hourly_logger_test\n")

	logger := slog.NewHourlyLoggerSt("houly_logger", "houly_logger.txt")

	for i := 0; i < 20; i++ {
		logger.Debug("houlry_logger test")
	}
}

// 你创建的一个logger，同一条日志可以根据需要打印到多个地方
// 比如stdout, stderr, file 之类的
func Muti_sink_test() {
	logger := slog.NewLogger("muti_logger")
	sink1 := slog.NewSimpleFileSinkSt("muti_sink_logger.txt")
	sink2 := slog.NewHourlyFileSinkSt("muti_hourly_logger.txt")
	sink3 := slog.NewStdoutSinkSt()

	logger.AppendSink(sink1).AppendSink(sink2).AppendSink(sink3)

	for i := 0; i < 20; i++ {
		logger.Debug("muti_sink_test test")
	}
}

func logger_log_test(gid int, logger *slog.Logger) {

	fmt.Printf("logger_log_test gid[%d]\n", gid)

	for i := 0; i < 5; i++ {
		logger.Debug("logger_log_test gid[%d] msgid[%d]", gid, i)
	}
}

func Muti_goroutine_stdout_test_nolock() {

	fmt.Printf("Muti_goroutine_stdout_test_nolock")

	logger := slog.NewStdoutLoggerSt("Muti_goroutine_test_nolock")

	for i := 0; i < 5; i++ {
		go logger_log_test(i, logger)
	}

	fmt.Printf("try sleep for a while\n")
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("sleep finished, Muti_goroutine_test_nolock end\n")
}

func Muti_goroutine_stdout_test_lock() {
	fmt.Printf("Muti_goroutine_stdout_test_lock")

	logger := slog.NewStdoutLoggerMt("Muti_goroutine_test_lock")

	for i := 0; i < 5; i++ {
		go logger_log_test(i, logger)
	}

	fmt.Printf("try sleep for a while\n")
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("sleep finished, Muti_goroutine_test_nolock end\n")
}

func Muti_goroutine_log_file_test_nolock() {

	fmt.Printf("Muti_goroutine_log_file_test_nolock")

	logger := slog.NewBasicLoggerSt("Muti_goroutine_log_file_test_nolock", "Muti_goroutine_log_file_test_nolock.txt")

	for i := 0; i < 5; i++ {
		go logger_log_test(i, logger)
	}

	fmt.Printf("try sleep for a while\n")
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("sleep finished, Muti_goroutine_test_nolock end\n")
}

func Muti_goroutine_log_file_test_lock() {
	fmt.Printf("Muti_goroutine_log_file_test_lock")

	logger := slog.NewBasicLoggerMt("Muti_goroutine_log_file_test_lock", "Muti_goroutine_log_file_test_lock.txt")

	for i := 0; i < 5; i++ {
		go logger_log_test(i, logger)
	}

	fmt.Printf("try sleep for a while\n")
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("sleep finished, Muti_goroutine_test_nolock end\n")
}
