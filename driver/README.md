DistLock 驱动程序的接口定义，共有3个接口，详见如下说明


## IWatcher

订阅接口，实现本接口的驱动可以监听锁的释放事件，及时唤醒等待中的 `Goroutine`

```go
type IWatcher interface {
	Watch(name string) (notifyChan <-chan struct{})
}
```


## IDriver

互斥锁接口，实现本接口的驱动即可用于互斥锁 `mutex` 的使用，返回值 `ok` 标识是否成功获取到锁，当为 `false` 时 `wait` 标识需要等待的时长，若驱动无法获取 `wait` 则返回负数即可使用默认等待时长 [defaultWait](https://github.com/go-locks/distlock/blob/master/mutex/options.go#L25)

```go
type IDriver interface {
	Lock(name, value string, expiry time.Duration) (ok bool, wait time.Duration)
	Unlock(name, value string)
	Touch(name, value string, expiry time.Duration) (ok bool)
}
```


## IRWDriver

读写锁接口，实现本接口的驱动即可用于读写锁 `rwmutex` 的使用，返回值意义同 `IDriver` 接口

```go
type IRWDriver interface {
	RLock(name, value string, expiry time.Duration) (ok bool, wait time.Duration)
	RUnlock(name, value string)
	RTouch(name, value string, expiry time.Duration) (ok bool)
	WLock(name, value string, expiry time.Duration) (ok bool, wait time.Duration)
	WUnlock(name, value string)
	WTouch(name, value string, expiry time.Duration) (ok bool)
}
```