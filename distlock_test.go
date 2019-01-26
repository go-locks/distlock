package distlock

import (
	"testing"

	"github.com/go-locks/distlock/driver/mocks"
	"github.com/stretchr/testify/mock"
)

var dl = mockDistlock()

func mockDistlock() *Distlock {
	var mockDriver = mocks.New()
	var notifyChan = make(chan struct{})
	mockDriver.On("Watch", mock.Anything).Return(notifyChan)
	return New(mockDriver, Prefix("test"))
}

func TestDistlock_NewMutex(t *testing.T) {
	name := "mutex"
	if mtx1, err := dl.NewMutex(name); err != nil {
		t.Error(err)
	} else if mtx2, _ := dl.NewMutex(name); mtx1 != mtx2 {
		t.Error("unexpected result, mutex must be same when the name is same")
	} else if _, err := dl.NewRWMutex(name); err == nil {
		t.Error("unexpected result, rwmutex can not have the same name with mutex")
	}
}

func TestDistlock_NewRWMutex(t *testing.T) {
	name := "rwmutex"
	if rwMtx1, err := dl.NewRWMutex(name); err != nil {
		t.Error(err)
	} else if rwMtx2, _ := dl.NewRWMutex(name); rwMtx1 != rwMtx2 {
		t.Error("unexpected result, rw mutex must be same when the name is same")
	} else if _, err := dl.NewMutex(name); err == nil {
		t.Error("unexpected result, mutex can not have the same name with rwmutex")
	}
}
