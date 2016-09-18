package gocoroutine

import (
	//	"fmt"
	"sync"
)

type taskSignal int

const (
	taskSignal_None taskSignal = iota
	taskSignal_Resume
)

type FlowControl interface {
	Yield(costFunc func(FlowControl))
	Params() interface{}
}

type Task struct {
	signalChan chan taskSignal

	executor func(FlowControl)
	params   interface{}

	needResumeGuard sync.Mutex
	needResume      bool

	sch *Scheduler
}

func (self *Task) Params() interface{} {
	return self.params
}

// 任务线程调用
func (self *Task) Yield(costFunc func(FlowControl)) {

	// 并行执行耗时任务
	go func() {

		costFunc(self)

		// 耗时任务完成
		self.needResumeGuard.Lock()
		self.needResume = true
		self.needResumeGuard.Unlock()

		self.sch.postCostTask(self)
	}()

	// 通知调度 本线程挂起
	self.sch.signal(scheduler_Yield)

	// 等待恢复逻辑
	<-self.signalChan
}

func (self *Task) NeedResume() bool {
	self.needResumeGuard.Lock()
	defer self.needResumeGuard.Unlock()
	return self.needResume
}

func (self *Task) signal(s taskSignal) {
	self.signalChan <- s
}

func NewTask() *Task {
	self := &Task{
		signalChan: make(chan taskSignal),
	}

	return self
}
