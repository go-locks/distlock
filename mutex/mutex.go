package mutex

import (
	"sync"

	"github.com/go-locks/distlock/driver"
)

type Mutex struct{ internal }

func NewMutex(name string, driver driver.IDriver, optFuncs ...OptFunc) *Mutex {
	locker := newLocker(name, driver, optFuncs...)
	return &Mutex{newInternal(locker, new(sync.Mutex))}
}

type RWMutex struct {
	readMtx  *Mutex
	writeMtx *Mutex
}

func NewRWMutex(name string, driver driver.IRWDriver, optFuncs ...OptFunc) *RWMutex {
	rwMutex := new(sync.RWMutex)
	readInternal := newInternal(newLocker(name, &readDriver{driver}, optFuncs...), &localReadMtx{rwMutex})
	writeInternal := newInternal(newLocker(name, &writeDriver{driver}, optFuncs...), rwMutex)
	return &RWMutex{readMtx: &Mutex{readInternal}, writeMtx: &Mutex{writeInternal}}
}

func (rwm *RWMutex) Read() *Mutex { return rwm.readMtx }

func (rwm *RWMutex) Write() *Mutex { return rwm.writeMtx }
