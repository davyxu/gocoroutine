package gocoroutine

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

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

func TestCoroutine(t *testing.T) {

	sch := NewScheduler()

	sch.Start()

	recvProc(sch)

	sch.Exit()

}
