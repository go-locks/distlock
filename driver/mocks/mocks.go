package mocks

import (
	"time"
	
	"github.com/stretchr/testify/mock"
)

type mocksDriver struct {
	mock.Mock
}

func New() *mocksDriver {
	return &mocksDriver{}
}

func (md *mocksDriver) Lock(name, value string, expiry time.Duration) (bool, time.Duration) {
	args := md.Called(name, value, expiry)
	return args.Bool(0), time.Duration(args.Int(1)) * time.Millisecond
}

func (md *mocksDriver) Unlock(name, value, channel string) {}

func (md *mocksDriver) Touch(name, value string, expiry time.Duration) bool {
	args := md.Called(name, value, expiry)
	return args.Bool(0)
}

func (md *mocksDriver) Watch(channel string) <-chan struct{} {
	args := md.Called(channel)
	return args.Get(0).(chan struct{})
}

func (md *mocksDriver) RLock(name, value string, expiry time.Duration) (bool, time.Duration) { return md.Lock(name, value, expiry) }

func (md *mocksDriver) RUnlock(name, value, channel string) {}

func (md *mocksDriver) RTouch(name, value string, expiry time.Duration) bool { return md.Touch(name, value, expiry) }

func (md *mocksDriver) WLock(name, value string, expiry time.Duration) (bool, time.Duration) { return md.Lock(name, value, expiry) }

func (md *mocksDriver) WUnlock(name, value, channel string) {}

func (md *mocksDriver) WTouch(name, value string, expiry time.Duration) bool { return md.Touch(name, value, expiry) }
