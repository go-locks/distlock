package mutex

import (
	"sync"
	"time"

	"github.com/go-locks/distlock/driver"
)

/* localReadMtx transform the read part of sync.RWMutex to sync.Locker
 * in this way, the internal just need to pay attention to ILocalMtx */

type localReadMtx struct {
	syncRWMutex *sync.RWMutex
}

func (lrMtx *localReadMtx) Lock() {
	lrMtx.syncRWMutex.RLock()
}

func (lrMtx *localReadMtx) Unlock() {
	lrMtx.syncRWMutex.RUnlock()
}

/* readDriver transform the read part of IRWDriver to IDriver,
 * in this way, the locker can just need to pay attention to IDriver */

type readDriver struct{ driver driver.IRWDriver }

func (rd *readDriver) Lock(name, value string, expiry time.Duration) (ok bool, wait time.Duration) {
	return rd.driver.RLock(name, value, expiry)
}

func (rd *readDriver) Unlock(name, value string) {
	rd.driver.RUnlock(name, value)
}

func (rd *readDriver) Touch(name, value string, expiry time.Duration) (ok bool) {
	return rd.driver.RTouch(name, value, expiry)
}

func (rd *readDriver) Watch(name string) <-chan struct{} {
	if wd, ok := rd.driver.(driver.IWatcher); ok {
		return wd.Watch(name)
	} else {
		return make(<-chan struct{})
	}
}

/* writeDriver transform the write part of IRWDriver to IDriver,
 * in this way, the locker can just need to pay attention to IDriver */

type writeDriver struct{ driver driver.IRWDriver }

func (wd *writeDriver) Lock(name, value string, expiry time.Duration) (ok bool, wait time.Duration) {
	return wd.driver.WLock(name, value, expiry)
}

func (wd *writeDriver) Unlock(name, value string) {
	wd.driver.WUnlock(name, value)
}

func (wd *writeDriver) Touch(name, value string, expiry time.Duration) (ok bool) {
	return wd.driver.WTouch(name, value, expiry)
}

func (wd *writeDriver) Watch(name string) <-chan struct{} {
	if wd, ok := wd.driver.(driver.IWatcher); ok {
		return wd.Watch(name)
	} else {
		return make(<-chan struct{})
	}
}
