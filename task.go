package gocoroutine

import (
	//	"fmt"
	"sync"
)

type TaskState int

const (
	TaskState_Running      TaskState = iota
	TaskState_Yield                  // 耗时任务同时开始
	TaskState_CostTaskDone           // 等待恢复
	TaskState_Resume                 // 已经恢复
	TaskState_WaitSync               // 等待调度同步
	TaskState_Done
)

const (
	TaskCommand_None = iota
	TaskCommand_Yield
	TaskCommand_Finished
)

type FlowControl interface {
	Yield(costFunc func(FlowControl))
	Params() interface{}
}

type Task struct {
	mailbox

	executor func(FlowControl)
	params   interface{}

	stateGuard sync.Mutex
	state      TaskState

	msgid int

	sch *Scheduler
}

func execWrapper(task *Task) {

	task.executor(task)

	task.Finished()

	task.sch.markFinished()
}

func costWrapper(task *Task, executor func(FlowControl)) {

	executor(task)

	task.CostTaskDone()
}

func (self *Task) Params() interface{} {
	return self.params
}

// 任务线程调用
func (self *Task) Yield(costFunc func(FlowControl)) {

	go costWrapper(self, costFunc)

	self.setState(TaskState_Yield)

	self.sch.send(TaskCommand_Yield)

	// 等待恢复逻辑
	self.recv()
}

// 任务线程调用
func (self *Task) Finished() {

	self.setState(TaskState_WaitSync)

	// 通知调度已经挂起
	self.sch.send(TaskCommand_Finished)

	self.setState(TaskState_Done)
}

// 耗时线程调用
func (self *Task) CostTaskDone() {

	self.setState(TaskState_CostTaskDone)
	self.sch.postTask(self, false)

}

func (self *Task) NeedResume() bool {
	return self.getState() == TaskState_CostTaskDone
}

func (self *Task) setState(s TaskState) {
	self.stateGuard.Lock()
	self.state = s
	//fmt.Printf("%v setstate %v\n", self.params, s)
	self.stateGuard.Unlock()
}

func (self *Task) getState() TaskState {
	self.stateGuard.Lock()
	defer self.stateGuard.Unlock()
	return self.state
}

func NewTask() *Task {
	self := &Task{

		mailbox: mailBoxAlloc(),
	}

	self.setState(TaskState_Running)

	return self
}
