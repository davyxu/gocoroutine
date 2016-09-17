package gocoroutine

import (
	//	"fmt"
	"sync"
)

const (
	SchedulerCommand_None = iota
	SchedulerCommand_Resume
)

type Scheduler struct {
	taskList chan *Task

	runningTask      int
	runningTaskGuard sync.Mutex

	mailbox
}

func (self *Scheduler) AddTask(exec func(FlowControl), params interface{}) {

	task := NewTask()
	task.executor = exec
	task.params = params
	task.sch = self

	self.postTask(task, true)
}

func (self *Scheduler) postTask(t *Task, normaltask bool) {

	if normaltask {
		self.runningTaskGuard.Lock()
		self.runningTask++
		self.runningTaskGuard.Unlock()

	}

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

	//	fmt.Println("markFinished", runningTaskLen())
}

func (self *Scheduler) fetchTask() *Task {

	return <-self.taskList
}

func (self *Scheduler) procTask() {

	for {

		task := self.fetchTask()

		if task == nil {
			break
		}

		if task.NeedResume() {

			// 通知任务线程完成
			task.send(SchedulerCommand_Resume)

		} else {
			go execWrapper(task)
		}

		self.mailbox.recv()
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
		taskList: make(chan *Task, 10),
		mailbox:  mailBoxAlloc(),
	}
}

func mailBoxAlloc() mailbox {
	return newChanBox()
}
