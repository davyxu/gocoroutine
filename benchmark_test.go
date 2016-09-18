package gocoroutine

import (
	//	"fmt"
	"sync"
	"testing"
)

func msgProc_benchmark(fc FlowControl) {

	//fmt.Println("msg 1")
	fc.Yield(dbProc_benchmark)

}

func dbProc_benchmark(fc FlowControl) {

	//fmt.Println("db 1")
	//	msgid := fc.Params().(int)
}

func recvProc_benchmark(sch *Scheduler) {

	for i := 0; i < 100000; i++ {
		sch.AddTask(msgProc_benchmark, i)
	}

}

func TestBenchmark(t *testing.T) {

	sch := NewScheduler()

	sch.Start()

	recvProc_benchmark(sch)

	sch.Exit()

}

func traditionalLogic(wg *sync.WaitGroup) {

	wg.Add(1)

	go func() {

		wg.Done()
	}()
}

// go test -bench=Benchmark -cpuprofile=cprof
func TestTraditional(t *testing.T) {

	var wg sync.WaitGroup

	for i := 0; i < 100000; i++ {

		traditionalLogic(&wg)
	}

	wg.Wait()

}
