package gocoroutine

import (
	"sync"
)

type mailbox interface {
	send(int)
	recv() int
}

type mutexBox struct {
	guard sync.Mutex
	con   *sync.Cond

	data      int
	dataGuard sync.Mutex
}

func (self *mutexBox) send(data int) {
	self.dataGuard.Lock()
	self.data = data
	self.dataGuard.Unlock()

	self.guard.Lock()
	self.con.Signal()
	self.guard.Unlock()
}

func (self *mutexBox) recv() int {
	self.guard.Lock()
	self.con.Wait()
	self.guard.Unlock()

	self.dataGuard.Lock()
	defer self.dataGuard.Unlock()

	return self.data

}

func newMutexBox() *mutexBox {
	self := &mutexBox{}

	self.con = sync.NewCond(&self.guard)

	return self
}

type channelBox struct {
	c chan int
}

func (self *channelBox) send(data int) {
	self.c <- data
}

func (self *channelBox) recv() int {
	return <-self.c
}

func newChanBox() *channelBox {
	return &channelBox{
		c: make(chan int),
	}
}
