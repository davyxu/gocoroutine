package gocoroutine

import (
	"sync"
	"testing"
)

func msgProc_benchmark(fc FlowControl) {

	fc.Yield(dbProc_benchmark)

}

func dbProc_benchmark(fc FlowControl) {

	//	msgid := fc.Params().(int)
}

func recvProc_benchmark() {

	for i := 0; i < 100000; i++ {

		go func(msgid int) {

			task := newTask()
			task.executor = msgProc_benchmark
			task.params = msgid

			//	fmt.Println(msgid, "recv Msg")
			postTask(task, true)

		}(i)

	}

}

func TestBenchmark(t *testing.T) {

	go procTask()
	recvProc_benchmark()

	waitExit()

}

func traditionalLogic(wg *sync.WaitGroup) {

	wg.Add(1)

	go func() {

		wg.Done()
	}()
}

func TestTraditional(t *testing.T) {

	var wg sync.WaitGroup

	for i := 0; i < 100000; i++ {

		traditionalLogic(&wg)
	}

	wg.Wait()

}
