package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type mocksDriver struct {
	mock.Mock
	until time.Time
}

func New() *mocksDriver {
	return &mocksDriver{}
}

func (md *mocksDriver) Lock(name, value string, expiry time.Duration) (bool, time.Duration) {
	args := md.Called(name, value, expiry)
	if args.Bool(0) {
		md.until = time.Now().Add(expiry)
	}
	return args.Bool(0), time.Duration(args.Int(1)) * time.Millisecond
}

func (md *mocksDriver) Unlock(name, value string) {
	md.until = time.Now()
}

func (md *mocksDriver) Touch(name, value string, expiry time.Duration) bool {
	if time.Now().Before(md.until) {
		md.until = time.Now().Add(expiry)
		return true
	}
	return false
}

func (md *mocksDriver) Watch(name string) <-chan struct{} {
	args := md.Called(name)
	return args.Get(0).(chan struct{})
}

func (md *mocksDriver) RLock(name, value string, expiry time.Duration) (bool, time.Duration) {
	return md.Lock(name, value, expiry)
}

func (md *mocksDriver) RUnlock(name, value string) {
	md.Unlock(name, value)
}

func (md *mocksDriver) RTouch(name, value string, expiry time.Duration) bool {
	return md.Touch(name, value, expiry)
}

func (md *mocksDriver) WLock(name, value string, expiry time.Duration) (bool, time.Duration) {
	return md.Lock(name, value, expiry)
}

func (md *mocksDriver) WUnlock(name, value string) {
	md.Unlock(name, value)
}

func (md *mocksDriver) WTouch(name, value string, expiry time.Duration) bool {
	return md.Touch(name, value, expiry)
}
