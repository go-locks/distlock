package mutex

import (
	"context"
	"time"

	"github.com/go-locks/distlock/driver"
)

type locker struct {
	options
	value      string
	until      time.Time
	driver     driver.IDriver
	notifyChan <-chan struct{}
}

func newLocker(name string, drve driver.IDriver, optFuncs ...OptFunc) locker {
	opts := newOptions(name, optFuncs...)
	locker := locker{options: opts, driver: drve}
	if wd, ok := locker.driver.(driver.IWatcher); ok {
		locker.notifyChan = wd.Watch(locker.name)
	} else {
		locker.notifyChan = make(<-chan struct{})
	}
	return locker
}

func (l *locker) lock() {
	value := generateValue()
	for {
		if ok, wait := l._lock(value); ok {
			return
		} else if wait > 0 {
			select {
			case <-l.notifyChan:
			case <-time.After(wait):
			}
		}
	}
}

func (l *locker) tryLock() bool {
	ok, _ := l._lock(generateValue())
	return ok
}

func (l *locker) lockCtx(ctx context.Context) bool {
	value := generateValue()
	for {
		if ok, wait := l._lock(value); ok {
			return true
		} else if wait > 0 {
			select {
			case <-ctx.Done():
				return false
			case <-l.notifyChan:
			case <-time.After(wait):
			}
		} else {
			select {
			case <-ctx.Done():
				return false
			default:
				continue
			}
		}
	}
}

func (l *locker) touch() bool {
	start := time.Now()
	if ok := l.driver.Touch(l.name, l.value, l.expiry); ok {
		cost := time.Now().Sub(start)
		left := l.costTopLimit - cost
		if left > 0 {
			l.until = time.Now().Add(left)
			return true
		}
	}
	return false
}

func (l *locker) unlock() {
	l.driver.Unlock(l.name, l.value)
}

func (l *locker) _lock(value string) (bool, time.Duration) {
	start := time.Now()
	ok, wait := l.driver.Lock(l.name, value, l.expiry)
	if ok {
		cost := time.Now().Sub(start)
		left := l.costTopLimit - cost
		if left > 0 {
			l.value = value
			l.until = time.Now().Add(left)
			return true, 0
		}
	} else if wait < 0 {
		wait = l.defaultWait
	}
	l.driver.Unlock(l.name, value)
	return false, wait
}
