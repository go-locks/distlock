package driver

import (
	"time"
)

type IWatcher interface {
	Watch(name string) (notifyChan <-chan struct{})
}

type IDriver interface {
	Lock(name, value string, expiry time.Duration) (ok bool, wait time.Duration)
	Unlock(name, value string)
	Touch(name, value string, expiry time.Duration) (ok bool)
}

type IRWDriver interface {
	RLock(name, value string, expiry time.Duration) (ok bool, wait time.Duration)
	RUnlock(name, value string)
	RTouch(name, value string, expiry time.Duration) (ok bool)
	WLock(name, value string, expiry time.Duration) (ok bool, wait time.Duration)
	WUnlock(name, value string)
	WTouch(name, value string, expiry time.Duration) (ok bool)
}
