package mutex

import (
	"context"
	"sync"
	"time"
)

type internal struct {
	locker   locker
	localMtx sync.Locker
}

func newInternal(locker locker, localMtx sync.Locker) internal {
	return internal{locker: locker, localMtx: localMtx}
}

func (itn *internal) Lock() {
	itn.localMtx.Lock()
	itn.locker.lock()
	itn.localMtx.Unlock()
}

func (itn *internal) LockCtx(ctx context.Context) bool {
	itn.localMtx.Lock()
	defer itn.localMtx.Unlock()
	return itn.locker.lockCtx(ctx)
}

func (itn *internal) TryLock() bool {
	itn.localMtx.Lock()
	defer itn.localMtx.Unlock()
	return itn.locker.tryLock()
}

func (itn *internal) Touch() bool {
	itn.localMtx.Lock()
	defer itn.localMtx.Unlock()
	return itn.locker.touch()
}

func (itn *internal) Unlock() {
	itn.localMtx.Lock()
	itn.locker.unlock()
	itn.localMtx.Unlock()
}

func (itn *internal) Until() time.Time {
	itn.localMtx.Lock()
	defer itn.localMtx.Unlock()
	return itn.locker.until
}
