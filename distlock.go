package distlock

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-locks/distlock/driver"
	"github.com/go-locks/distlock/mutex"
)

type Distlock struct {
	sync.Mutex
	prefix   string
	driver   driver.IDriver
	mtxMap   map[string]*mutex.Mutex
	rwMtxMap map[string]*mutex.RWMutex
}

func New(driver driver.IDriver, optFuncs ...OptFunc) *Distlock {
	dl := Distlock{
		driver:   driver,
		prefix:   "distlock-",
		mtxMap:   make(map[string]*mutex.Mutex),
		rwMtxMap: make(map[string]*mutex.RWMutex),
	}
	for _, optFunc := range optFuncs {
		optFunc(&dl)
	}
	return &dl
}

func (dl *Distlock) buildName(name string) string {
	return dl.prefix + name
}

func (dl *Distlock) NewMutex(name string, optFuncs ...mutex.OptFunc) (*mutex.Mutex, error) {
	dl.Lock()
	defer dl.Unlock()
	if _, ok := dl.mtxMap[name]; ok {
		return dl.mtxMap[name], nil
	} else if _, ok := dl.rwMtxMap[name]; ok {
		return nil, fmt.Errorf("a rw mutex named '%s' already exist, instead of mutex", name)
	}
	dl.mtxMap[name] = mutex.NewMutex(dl.buildName(name), dl.driver, optFuncs...)
	return dl.mtxMap[name], nil
}

func (dl *Distlock) NewRWMutex(name string, optFuncs ...mutex.OptFunc) (*mutex.RWMutex, error) {
	dl.Lock()
	defer dl.Unlock()
	rwDriver, ok := dl.driver.(driver.IRWDriver)
	if !ok {
		return nil, errors.New("the driver is not a rw driver, so you can not create rw mutex")
	}
	if _, ok := dl.rwMtxMap[name]; ok {
		return dl.rwMtxMap[name], nil
	} else if _, ok := dl.mtxMap[name]; ok {
		return nil, fmt.Errorf("a mutex named '%s' already exist, instead of rw mutex", name)
	}
	dl.rwMtxMap[name] = mutex.NewRWMutex(dl.buildName(name), rwDriver, optFuncs...)
	return dl.rwMtxMap[name], nil
}
