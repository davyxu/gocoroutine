# gocoroutine
为游戏服务器设计的coroutine实现

# 传统逻辑编写法-阻塞同步多线程

遇到db请求时, 阻塞等待, db完成时继续逻辑

多个连接使用线程池或者goroutine进行编写

## 优点

* 逻辑清晰

* 适用于玩家间交互不多的情况, http服务器就是这种模型

## 缺点

* 遇到公共逻辑, 诸如公会, 排行榜之类共享数据时, 进出均需要加锁

* 对开发人员要求高, 必须有多线程安全意识


# 传统逻辑编写法-非阻塞异步单线程

遇到db请求时, 将请求在db线程进行处理, 通过回调通知逻辑线程

## 优点

* 单线程逻辑,对开发人员编写逻辑要求不高

* 适用于玩家间交互多, 且保证逻辑裁决单一的情况, MMORPG服务器就是这种模型

## 缺点

* 在多个异步请求迭代时, 容易陷入callback-hell. 逻辑跳跃,需要配合文档理解编写者思想


# gocoroutine编写法-非阻塞同步单线程

* 将每一个消息处理划分为1个任务

* 任务由调度器管理, 调度器负责与任务进行交互, 控制各种异步, 同步流程

* 消息线程负责将消息上下文创建为任务, 并投递到调度器的任务队列

* 调度器从任务队列获取任务, 将任务放在独立的goroutine里

* 遇到db请求时, 将请求在db线程进行处理, 挂起当前任务的进度, 通知调度器取下一个任务

* 如果没有遇到db需求直到任务完成, 任务通知调度器取下一个任务

* 当db耗时任务完成时, 将任务重新投放到任务队列中, 并标记可唤醒

* 调度器取出需要唤醒的任务时, 通知任务取消挂起继续执行直到完成

## 优点

* 逻辑清晰

* 单线程逻辑,对开发人员编写逻辑要求不高

## 缺点

* 相比非阻塞异步单线程, 在性能上有少许损耗

```golang
// 消息处理
func msgProc(fc FlowControl) {

	fc.Yield(dbProc)

}

// db耗时处理
func dbProc(fc FlowControl) {

	msgid := fc.Params().(int)

	du := rand.Int31n(5) + 1
	time.Sleep(time.Duration(du) * time.Second)
	fmt.Println(msgid, "dbProc")
}

// 消息线程, 生成任务
func recvProc(sch *Scheduler) {

	for i := 0; i < 2; i++ {

		sch.AddTask(msgProc, i)

		fmt.Println(i, "recv Msg")

	}

}

func main( ) {

	sch := NewScheduler()

	sch.Start()

	recvProc(sch)

	sch.Exit()

}

```


# 备注

感觉不错请star, 谢谢!

知乎: [http://www.zhihu.com/people/sunicdavy](http://www.zhihu.com/people/sunicdavy)
