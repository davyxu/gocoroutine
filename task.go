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

type FlowControl interface {
	Yield(costFunc func(FlowControl))
	Params() interface{}
}

type Task struct {
	scheuler2task chan string
	task2scheuler chan string

	executor func(FlowControl)
	params   interface{}

	stateGuard sync.Mutex
	state      TaskState

	msgid int
}

func execWrapper(task *Task) {

	task.executor(task)

	task.Finished()

	markFinished()
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

	addCostTask(self)
	go costWrapper(self, costFunc)

	self.setState(TaskState_Yield)

	// 通知调度已经挂起
	self.task2scheuler <- "yield"

	// 等待恢复逻辑
	<-self.scheuler2task

}

// 任务线程调用
func (self *Task) Finished() {

	self.setState(TaskState_WaitSync)

	// 通知调度已经挂起
	self.task2scheuler <- "finished"

	self.setState(TaskState_Done)
}

// 耗时线程调用
func (self *Task) CostTaskDone() {

	self.setState(TaskState_CostTaskDone)
	postTask(self, false)

	removeCostTask(self)
}

// 调度线程调用
func (self *Task) Exec() string {

	go execWrapper(self)

	result := <-self.task2scheuler

	//	fmt.Println(self.params, "scheuler", result)

	return result
}

func (self *Task) NeedResume() bool {
	return self.getState() == TaskState_CostTaskDone
}

// 调度线程调用
func (self *Task) Resume() string {

	self.setState(TaskState_Resume)

	// 通知任务线程完成
	self.scheuler2task <- "resume"

	result := <-self.task2scheuler

	//	fmt.Println(self.params, "scheuler", result)

	return result
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

func newTask() *Task {
	self := &Task{
		scheuler2task: make(chan string),
		task2scheuler: make(chan string),
	}

	self.setState(TaskState_Running)

	return self
}
