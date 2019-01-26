package mutex

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-locks/distlock/driver"
	"github.com/go-locks/distlock/driver/mocks"
	"github.com/stretchr/testify/mock"
)

var name = "test"
var expiry = 10
var factor = 0.10
var timeUnit = time.Millisecond
var notifyChan = make(chan struct{})
var watchChannel = WatchChannelPrefix + name
var expiryDuration = timeUnit * time.Duration(expiry)
var validDelayDuration = timeUnit * time.Duration(float64(expiry)*(1-factor)*50) / 100
var invalidDelayDuration = timeUnit * time.Duration(float64(expiry)*(1-factor)*100) / 100

func mockDriver(ok bool, wait int, delay time.Duration) interface{} {
	var mockDriver = mocks.New()
	mockDriver.On("Watch", watchChannel).Return(notifyChan)
	mockDriver.On("Lock", name, mock.Anything, expiryDuration).After(delay).Return(ok, wait)
	mockDriver.On("RLock", name, mock.Anything, expiryDuration).After(delay).Return(ok, wait)
	mockDriver.On("WLock", name, mock.Anything, expiryDuration).After(delay).Return(ok, wait)
	mockDriver.On("Touch", name, mock.Anything, expiryDuration).After(delay).Return(ok, wait)
	return mockDriver
}

func mockSuccessMutex() *Mutex {
	drve := mockDriver(true, 0, 0).(driver.IDriver)
	return NewMutex(name, drve, Expiry(expiryDuration), Factor(factor))
}

func mockFailureMutex() *Mutex {
	drve := mockDriver(false, 100, 0).(driver.IDriver)
	return NewMutex(name, drve, Expiry(expiryDuration), Factor(factor))
}

func mockDelaySuccessMutex() *Mutex {
	drve := mockDriver(true, 0, validDelayDuration).(driver.IDriver)
	return NewMutex(name, drve, Expiry(expiryDuration), Factor(factor))
}

func mockDelayFailureMutex() *Mutex {
	drve := mockDriver(true, 0, invalidDelayDuration).(driver.IDriver)
	return NewMutex(name, drve, Expiry(expiryDuration), Factor(factor))
}

func mockSuccessRWMutex() *RWMutex {
	drve := mockDriver(true, 0, 0).(driver.IRWDriver)
	return NewRWMutex(name, drve, Expiry(expiryDuration), Factor(factor))
}

func TestMutex_Lock(t *testing.T) {
	var counter int32
	go func() {
		mockSuccessMutex().Lock()
		atomic.AddInt32(&counter, 1) // +1
	}()
	go func() {
		mockDelaySuccessMutex().Lock()
		atomic.AddInt32(&counter, 1) // +1
	}()
	go func() {
		mockDelayFailureMutex().Lock()
		atomic.AddInt32(&counter, 1)
	}()
	go func() {
		mockFailureMutex().Lock()
		atomic.AddInt32(&counter, 1)
	}()
	time.Sleep(expiryDuration)
	if counter != 2 {
		t.Errorf("unexpected result, expect = %d, but = %d", 2, counter)
	}
}

func TestMutex_LockCtx(t *testing.T) {
	var counter int32
	go func() {
		if mockSuccessMutex().LockCtx(context.TODO()) {
			atomic.AddInt32(&counter, 1) // +1
		}
	}()
	go func() {
		var ctx, _ = context.WithTimeout(context.TODO(), expiryDuration-time.Millisecond)
		if mockFailureMutex().LockCtx(ctx) {
			atomic.AddInt32(&counter, 1)
		}
		atomic.AddInt32(&counter, 1) // +1
	}()
	go func() {
		var ctx, _ = context.WithTimeout(context.TODO(), expiryDuration+time.Millisecond)
		if mockFailureMutex().LockCtx(ctx) {
			atomic.AddInt32(&counter, 1)
		}
		atomic.AddInt32(&counter, 1)
	}()
	time.Sleep(expiryDuration)
	if counter != 2 {
		t.Errorf("unexpected result, expect = %d, but = %d", 2, counter)
	}
}

func TestMutex_TryLock(t *testing.T) {
	if !mockSuccessMutex().TryLock() {
		t.Error("unexpected result, expect = true, but = false")
	}
	if mockFailureMutex().TryLock() {
		t.Error("unexpected result, expect = false, but = true")
	}
}

func TestMutex_Touch(t *testing.T) {
	if !mockSuccessMutex().Touch() {
		t.Error("unexpected result, expect = true, but = false")
	}
	if mockFailureMutex().Touch() {
		t.Error("unexpected result, expect = false, but = true")
	}
}

func TestRWMutex_ReadLock(t *testing.T) {
	var counter int32
	go func() {
		mockSuccessRWMutex().Read().Lock()
		atomic.AddInt32(&counter, 1) // + 1
	}()
	go func() {
		if mockSuccessRWMutex().Read().LockCtx(context.TODO()) {
			atomic.AddInt32(&counter, 1) // + 1
		}
		if mockSuccessRWMutex().Read().TryLock() {
			atomic.AddInt32(&counter, 1) // + 1
		}
	}()
	time.Sleep(expiryDuration)
	if counter != 3 {
		t.Errorf("unexpected result, expect = %d, but = %d", 3, counter)
	}
}
func TestRWMutex_WriteLock(t *testing.T) {
	var counter int32
	go func() {
		mockSuccessRWMutex().Write().Lock()
		atomic.AddInt32(&counter, 1) // + 1
	}()
	go func() {
		if mockSuccessRWMutex().Write().LockCtx(context.TODO()) {
			atomic.AddInt32(&counter, 1) // + 1
		}
		if mockSuccessRWMutex().Write().TryLock() {
			atomic.AddInt32(&counter, 1) // + 1
		}
	}()
	time.Sleep(expiryDuration)
	if counter != 3 {
		t.Errorf("unexpected result, expect = %d, but = %d", 3, counter)
	}
}
