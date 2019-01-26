package driver

import (
	"time"
)

type IWatcher interface {
	Watch(channel string) (notifyChan <-chan struct{})
}

type IDriver interface {
	Lock(name, value string, expiry time.Duration) (ok bool, wait time.Duration)
	Unlock(name, value, channel string)
	Touch(name, value string, expiry time.Duration) (ok bool)
}

type IRWDriver interface {
	RLock(name, value string, expiry time.Duration) (ok bool, wait time.Duration)
	RUnlock(name, value, channel string)
	RTouch(name, value string, expiry time.Duration) (ok bool)
	WLock(name, value string, expiry time.Duration) (ok bool, wait time.Duration)
	WUnlock(name, value, channel string)
	WTouch(name, value string, expiry time.Duration) (ok bool)
}
