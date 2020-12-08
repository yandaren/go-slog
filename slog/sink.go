// File sink.go
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

type Sink interface {
	// log the msg
	Log(msg string)

	// flush the log
	Flush()

	// check need log the message
	ShoudLog(lvl LogLevel) bool

	// set log level
	SetLogLvl(lvl LogLevel)

	// get log level
	GetLogLvl() LogLevel

	// set force flush
	SetForceFlush(force bool)

	// is force flush
	IsForceFlush() bool

	// lock
	Lock()

	// unlock
	Unlock()
}

type BaseSink struct {
	log_lvl     LogLevel // logger level
	force_flush bool     // force flush eveny log
	locker      Locker
}

func NewBaseSinkSt() *BaseSink {
	return &BaseSink{
		log_lvl:     LvlDebug,
		force_flush: false,
		locker:      NewNullLocker(),
	}
}

func NewBaseSinkMt() *BaseSink {
	return &BaseSink{
		log_lvl:     LvlDebug,
		force_flush: false,
		locker:      NewExclusiveLocker(),
	}
}

func (this *BaseSink) SetLocker(lk Locker) {
	this.locker = lk
}

func (this *BaseSink) Lock() {
	if this.locker != nil {
		this.locker.Lock()
	}
}

func (this *BaseSink) Unlock() {
	if this.locker != nil {
		this.locker.Unlock()
	}
}

func (this *BaseSink) DefaultInit() {
	this.log_lvl = LvlDebug
	this.force_flush = false
}

func (this *BaseSink) Log(msg string) {
}

func (this *BaseSink) Flush() {
}

func (this *BaseSink) ShoudLog(lvl LogLevel) bool {
	return lvl >= this.log_lvl
}

func (this *BaseSink) SetLogLvl(lvl LogLevel) {
	this.log_lvl = lvl
}

func (this *BaseSink) GetLogLvl() LogLevel {
	return this.log_lvl
}

func (this *BaseSink) SetForceFlush(force bool) {
	this.force_flush = force
}

func (this *BaseSink) IsForceFlush() bool {
	return this.force_flush
}
