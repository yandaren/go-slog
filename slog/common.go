// File common
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

import (
	"fmt"
	"strings"
)

// Level type
type LogLevel uint32

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	LvlDebug LogLevel = iota
	LvlInfo
	LvlWarn
	LvlError
	LvlFatal
	LvlNone
)

var AllLogLevels = []LogLevel{
	LvlDebug,
	LvlInfo,
	LvlWarn,
	LvlError,
	LvlFatal,
	LvlNone,
}

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level LogLevel) String() string {
	switch level {
	case LvlDebug:
		return "debug"
	case LvlInfo:
		return "info"
	case LvlWarn:
		return "warn"
	case LvlError:
		return "error"
	case LvlFatal:
		return "fatal"
	case LvlNone:
		return "none"
	}

	return "unknown"
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (LogLevel, error) {
	switch strings.ToLower(lvl) {
	case "none":
		return LvlNone, nil
	case "fatal":
		return LvlFatal, nil
	case "error":
		return LvlError, nil
	case "warn", "warning":
		return LvlWarn, nil
	case "info":
		return LvlInfo, nil
	case "debug":
		return LvlDebug, nil
	}

	var l LogLevel
	return l, fmt.Errorf("not a valid slog Level: %q", lvl)
}
