// File locker
// @Author: yandaren1220@126.com
// @Date: 2018-08-13

package slog

import (
	"sync"
)

// locker interface
type Locker interface {
	Lock()
	Unlock()
}

// exclusive lock
type ExclusiveLocker struct {
	mtx sync.Mutex // exclusive mtx
}

// create exclusive lock
func NewExclusiveLocker() *ExclusiveLocker {
	return &ExclusiveLocker{}
}

// override the unlock/unlock method
func (this *ExclusiveLocker) Lock() {
	this.mtx.Lock()
}

func (this *ExclusiveLocker) Unlock() {
	this.mtx.Unlock()
}

// readwrite lock
type ReadWriteLocker struct {
	mtx sync.RWMutex // rw mutx
}

// create ReadWrite locker
func NewReadWriteLocker() *ReadWriteLocker {
	return &ReadWriteLocker{}
}

// override the unlock/unlock method
func (this *ReadWriteLocker) Lock() {
	this.mtx.Lock()
}

func (this *ReadWriteLocker) Unlock() {
	this.mtx.Unlock()
}

// null locker
type NullLocker struct {
}

// create null locker
func NewNullLocker() *NullLocker {
	return &NullLocker{}
}

// override the unlock/unlock method
func (this *NullLocker) Lock() {
	// do nothing
}

func (this *NullLocker) Unlock() {
	// do nothing
}
