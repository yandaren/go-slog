// File slog.go
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

// go-slog project go-slog.go
package slog

// apis

func new_basic_logger(logger_name string, file_name string, thread_safe bool) *Logger {
	logger := NewLogger(logger_name)
	sink := NewSimpleFileSink(file_name, thread_safe)
	logger.AppendSink(sink)

	return logger
}

// create not thread safe basic simple file logger
func NewBasicLoggerSt(logger_name string, file_name string) *Logger {
	return new_basic_logger(logger_name, file_name, false)
}

// create thread safe basic simple file logger
func NewBasicLoggerMt(logger_name string, file_name string) *Logger {
	return new_basic_logger(logger_name, file_name, true)
}

func new_rotating_logger(logger_name string, base_file_name string, max_file_size int64, max_files int32, thread_safe bool) *Logger {
	logger := NewLogger(logger_name)
	sink := NewRotatingFileSink(base_file_name, max_file_size, max_files, thread_safe)
	logger.AppendSink(sink)

	return logger
}

// create not thread safe rotating logger
func NewRotatingLoggerSt(logger_name string, base_file_name string, max_file_size int64, max_files int32) *Logger {
	return new_rotating_logger(logger_name, base_file_name, max_file_size, max_files, false)
}

// create thread safe rotating logger
func NewRotatingLoggerMt(logger_name string, base_file_name string, max_file_size int64, max_files int32) *Logger {
	return new_rotating_logger(logger_name, base_file_name, max_file_size, max_files, true)
}

func new_daily_logger(logger_name string, base_file_name string, rotation_hour int32, rotation_minute int32, thread_safe bool) *Logger {
	sink, err := NewDailyFileSink(base_file_name, rotation_hour, rotation_minute, thread_safe)
	if err != nil {
		return nil
	}

	logger := NewLogger(logger_name)

	logger.AppendSink(sink)

	return logger
}

// create not thread safe daily logger
func NewDailyLoggerSt(logger_name string, base_file_name string, rotation_hour int32, rotation_minute int32) *Logger {
	return new_daily_logger(logger_name, base_file_name, rotation_hour, rotation_minute, false)
}

// create thread safe daily logger
func NewDailyLoggerMt(logger_name string, base_file_name string, rotation_hour int32, rotation_minute int32) *Logger {
	return new_daily_logger(logger_name, base_file_name, rotation_hour, rotation_minute, true)
}

func new_hourly_logger(logger_name string, base_file_name string, thread_safe bool) *Logger {
	logger := NewLogger(logger_name)
	sink := NewHourlyFileSink(base_file_name, thread_safe)
	logger.AppendSink(sink)

	return logger
}

// create not thread safe hourly file logger
func NewHourlyLoggerSt(logger_name string, base_file_name string) *Logger {
	return new_hourly_logger(logger_name, base_file_name, false)
}

// create thread safe hourly file logger
func NewHourlyLoggerMt(logger_name string, base_file_name string) *Logger {
	return new_hourly_logger(logger_name, base_file_name, true)
}

func new_stdout_logger(name string, thread_safe bool) *Logger {
	logger := NewLogger(name)
	sink := NewStdoutSink(thread_safe)
	logger.AppendSink(sink)

	return logger
}

// create not thread safe stdout logger
func NewStdoutLoggerSt(name string) *Logger {
	return new_stdout_logger(name, false)
}

// create thread safe  stdout logger
func NewStdoutLoggerMt(name string) *Logger {
	return new_stdout_logger(name, true)
}

func new_stderr_logger(name string, thread_safe bool) *Logger {
	logger := NewLogger(name)
	sink := NewStderrSink(thread_safe)
	logger.AppendSink(sink)

	return logger
}

// create not thread safe stderr logger
func NewStderrLoggerSt(name string) *Logger {
	return new_stderr_logger(name, false)
}

// create thread safe stderror logger
func NewStderrLoggerMt(name string) *Logger {
	return new_stderr_logger(name, true)
}
