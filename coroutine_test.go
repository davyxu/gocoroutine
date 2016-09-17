package gocoroutine

import (
	"fmt"
	"math/rand"
	///	"testing"
	"time"
)

func msgProc(fc FlowControl) {

	fc.Yield(dbProc)

}

func dbProc(fc FlowControl) {

	msgid := fc.Params().(int)

	du := rand.Int31n(5) + 1
	time.Sleep(time.Duration(du) * time.Second)
	fmt.Println(msgid, "dbProc")
}

func recvProc(sch *Scheduler) {

	for i := 0; i < 1; i++ {

		go func(msgid int) {

			sch.AddTask(msgProc_benchmark, msgid)

			fmt.Println(msgid, "recv Msg")

		}(i)

	}

}

//func TestCoroutine(t *testing.T) {

//	sch := NewScheduler()

//	sch.Start()

//	recvProc(sch)

//	sch.Exit()

//}
