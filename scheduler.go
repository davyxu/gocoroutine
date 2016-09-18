package gocoroutine

import (
	//	"fmt"
	"sync"
)

type schedulerSignal int

const (
	scheduler_None schedulerSignal = iota
	scheduler_Yield
	scheduler_Finished
)

type Scheduler struct {
	taskList chan *Task

	runningTask      int
	runningTaskGuard sync.Mutex

	signalChan chan schedulerSignal
}

func (self *Scheduler) AddTask(exec func(FlowControl), params interface{}) {

	task := NewTask()
	task.executor = exec
	task.params = params
	task.sch = self

	self.runningTaskGuard.Lock()
	self.runningTask++
	self.runningTaskGuard.Unlock()

	self.taskList <- task
}

func (self *Scheduler) postCostTask(t *Task) {
	self.taskList <- t
}

func (self *Scheduler) runningTaskLen() int {
	self.runningTaskGuard.Lock()

	defer self.runningTaskGuard.Unlock()

	return self.runningTask

}

func (self *Scheduler) markFinished() {
	self.runningTaskGuard.Lock()
	self.runningTask--
	self.runningTaskGuard.Unlock()
}

func (self *Scheduler) fetchTask() *Task {

	return <-self.taskList
}

func (self *Scheduler) signal(s schedulerSignal) {
	self.signalChan <- s
}

func (self *Scheduler) procTask() {

	for {

		task := self.fetchTask()

		if task == nil {
			break
		}

		if task.NeedResume() {

			// 通知任务线程恢复逻辑
			task.signal(taskSignal_Resume)

		} else {
			go func() {
				task.executor(task)

				// 通知完成
				self.signal(scheduler_Finished)

				self.markFinished()
			}()
		}

		// 等待完成或者挂起
		<-self.signalChan
	}

}

func (self *Scheduler) Start() {
	go self.procTask()
}

func (self *Scheduler) Exit() {

	endChan := make(chan bool)

	go func() {

		for {

			l := self.runningTaskLen()

			if l == 0 {
				break
			}

		}

		endChan <- true

	}()

	<-endChan

}

func NewScheduler() *Scheduler {

	return &Scheduler{
		taskList:   make(chan *Task, 10),
		signalChan: make(chan schedulerSignal),
	}
}
