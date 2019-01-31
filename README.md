# Go-Locks [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)

通用的Golang分布式锁组件，更多使用案例详见 [examples](https://github.com/go-locks/examples)


## Driver列表

若有意向贡献未完成的驱动代码，请通过 [ISSUES](https://github.com/go-locks/distlock/issues) 或 邮箱 `249008728@qq.com` 联系我

| Driver | 代码完成度 | 测试完成度 | 依赖包                                                | 使用说明                                                    |
| :----- | :--------: | :--------: | :---------------------------------------------------- | :---------------------------------------------------------- |
| redis  | 100%       | 100%       | [letsfire/redigo](https://github.com/letsfire/redigo) | 详见 [README.md](https://github.com/go-locks/redis-driver)  |
| pgsql  | 100%       | 100%       | [lib/pq](https://github.com/lib/pq)                   | 详见 [README.md](https://github.com/go-locks/pgsql-driver)  |
| etcd   | 未完成     | 未测试     | [etcd/client](https://go.etcd.io/etcd/client)         | 详见 [README.md](https://github.com/go-locks/etcd-driver)   |
| etcdv3 | 未完成     | 未测试     | [etcd/clientv3](https://go.etcd.io/etcd/clientv3)     | 详见 [README.md](https://github.com/go-locks/etcdv3-driver) |


## 方法说明

配置项 `mutex.OptFunc` 以及返回值锁的使用详见 [mutex/README.md](https://github.com/go-locks/distlock/tree/master/mutex)

- **NewMutex(name string, optFuncs ...mutex.OptFunc) (\*mutex.Mutex, error)**  
创建互斥锁，若 `name` 已用于创建读写锁则返回 `error`，本方法单例模式
- **NewRWMutex(name string, optFuncs ...mutex.OptFunc) (\*mutex.RWMutex, error)**  
创建读写锁，若 `name` 已用于创建互斥锁则返回 `error`，本方法单例模式


## 注意事项

* 不可重入（如果您有强烈的需求场景，请通过 [ISSUES](https://github.com/go-locks/distlock/issues) 提供反馈）
* 非公平锁（Golang的本地锁 `sync.Locker` 视乎也不是公平锁，若您有需求或建议，请通过 [ISSUES](https://github.com/go-locks/distlock/issues) 提供反馈）
* 有互斥锁 `mutex` 和 读写锁 `rwmutex` 两种类型，具体支持程度详见各个 `Driver` 对应的 `README.md`
* 虽有完整的单元测试，但暂未经过实际项目考验，故慎用于生产环境，如有问题请通过 [ISSUES](https://github.com/go-locks/distlock/issues) 来共同完善


## 项目结构

* 主线调用层级为 `distlock.go` -> `mutex.go` -> `driver.go`
* `distlock.go` 提供了创建锁的工厂类，单例模式（相同名称的锁有且仅有一个，有且仅为一种）
* `mutex.go`提供了各类锁的实现，欢迎各位同学贡献其他类型锁，详见 [mutex/README.md](https://github.com/go-locks/distlock/tree/master/mutex)
* `driver.go`提供驱动接口的定义，欢迎各位同学贡献其他驱动，详见 [driver/README.md](https://github.com/go-locks/distlock/tree/master/driver)
