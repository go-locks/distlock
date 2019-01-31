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
var expiryDuration = timeUnit * time.Duration(expiry)
var validDelayDuration = timeUnit * time.Duration(float64(expiry)*(1-factor)*50) / 100
var invalidDelayDuration = timeUnit * time.Duration(float64(expiry)*(1-factor)*100) / 100

func mockDriver(ok bool, wait int, delay time.Duration) interface{} {
	var mockDriver = mocks.New()
	mockDriver.On("Watch", name).Return(notifyChan)
	mockDriver.On("Lock", name, mock.Anything, expiryDuration).After(delay).Return(ok, wait)
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
	var addIfLocked = func(mtx *Mutex) {
		mtx.Lock()
		atomic.AddInt32(&counter, 1)
	}

	go addIfLocked(mockSuccessMutex())      // +1
	go addIfLocked(mockDelaySuccessMutex()) // +1
	go addIfLocked(mockFailureMutex())
	go addIfLocked(mockDelayFailureMutex())

	time.Sleep(expiryDuration)

	if counter != 2 {
		t.Errorf("unexpected result, expect = 2, but = %d", counter)
	}
}

func TestMutex_LockCtx(t *testing.T) {
	var counter int32
	var addIfLocked = func(mtx *Mutex, ctx context.Context) {
		if mtx.LockCtx(ctx) {
			atomic.AddInt32(&counter, 1)
		}
		atomic.AddInt32(&counter, 1)
	}
	var ctx1, _ = context.WithTimeout(context.TODO(), expiryDuration-time.Millisecond)
	var ctx2, _ = context.WithTimeout(context.TODO(), expiryDuration+time.Millisecond)

	go addIfLocked(mockSuccessMutex(), context.TODO()) // +2
	go addIfLocked(mockFailureMutex(), ctx1)           // +1
	go addIfLocked(mockFailureMutex(), ctx2)

	time.Sleep(expiryDuration)

	if counter != 3 {
		t.Errorf("unexpected result, expect = 3, but = %d", counter)
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
	mtx := mockSuccessMutex()
	if !mtx.TryLock() || !mtx.Touch() {
		t.Error("unexpected result, expect = true, but = false")
	}
	time.Sleep(expiryDuration)
	if mtx.Touch() {
		t.Error("unexpected result, expect = false, but = true")
	}
}

func TestMutex_Heartbeat(t *testing.T) {
	mtx := mockSuccessMutex()
	ctx, cancel := context.WithCancel(context.TODO())

	mtx.TryLock()
	mtx.Heartbeat(ctx)

	time.Sleep(expiryDuration)
	if !mtx.Touch() {
		t.Error("unexpected result, expect = true, but = false")
	}

	cancel()
	time.Sleep(expiryDuration)
	if mtx.Touch() {
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
