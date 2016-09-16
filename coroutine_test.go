package gocoroutine

import (
	"fmt"
	"math/rand"
	//	"testing"
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

func recvProc() {

	for i := 0; i < 10; i++ {

		go func(msgid int) {

			task := newTask()
			task.executor = msgProc
			task.params = msgid

			fmt.Println(msgid, "recv Msg")
			postTask(task, true)

		}(i)

	}

}

//func TestCoroutine(t *testing.T) {

//	go procTask()
//	recvProc()

//	waitExit()

//}
