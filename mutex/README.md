在锁驱动 [Driver](https://github.com/go-locks/distlock/tree/master/driver) 的基础提供封装，暴露外部可使用的方法


## 参数配置

- **expiry** - 默认 5s，锁的存活时长，根据使用场景设置合适的存活时长，避免锁释放失败后造成死锁
- **factor** - 默认 0.12，(factor * expiry)做为默认 `wait` 时长，也做为 `Until()` 的预留时长，详见 [options.go](options.go#L24-L26)


## 方法说明

- **Lock()**  
阻塞锁，永久阻塞直到成功获取锁
- **LockCtx(ctx context.Context) bool**  
带 `context` 的阻塞锁，永久阻塞直到成功获取锁（返回 `true`）或 外部通过 `context` 中断（返回 `false`）
- **TryLock() bool**  
非阻塞，尝试加锁，成功返回 `true`，失败立即返回 `false`
- **Touch() bool**  
在加锁成功后用于探测锁是否还持有，若持有则更新到期时间并返回 `true`，否则返回 `false`
- **Unlock()**  
释放锁，锁使用后应及时释放，若没有释放则只有等待锁到期后其他 `Goroutine` 才能获取
- **Until() time.Time**  
返回当前锁的到期时间点，若在到达该点时任务没有完成则必须通过 `Touch` 方法进行检测并延期
- **Heartbeat(ctx context.Context) <-chan struct{}**  
自带的锁续约保活机制，参数 `ctx` 用于退出心跳循环，返回值用于向外传递锁 `Touch` 失败的信号，若 `Touch` 失败心跳循环同样会中止。更复杂的情况请使用 `Until` 和 `Touch` 自行实现延期保活
