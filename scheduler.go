package gocoroutine

import (
	"sync"
)

var taskList = make(chan *Task, 10)

var runningTask int
var runningTaskGuard sync.Mutex

func postTask(t *Task, normaltask bool) {

	if normaltask {
		runningTaskGuard.Lock()
		runningTask++
		runningTaskGuard.Unlock()
	}

	taskList <- t
}

func runningTaskLen() int {
	runningTaskGuard.Lock()

	defer runningTaskGuard.Unlock()

	return runningTask

}

func markFinished() {

	runningTaskGuard.Lock()
	runningTask--
	runningTaskGuard.Unlock()

	//	fmt.Println("markFinished", runningTaskLen())
}

func fetchTask() *Task {

	return <-taskList
}

func procTask() {

	for {

		task := fetchTask()

		if task == nil {
			break
		}

		if task.NeedResume() {
			task.Resume()
		} else {

			task.Exec()
		}
	}

}

var costTaskMap = make(map[*Task]*Task)
var costTaskMapGuard sync.Mutex

func addCostTask(t *Task) {
	costTaskMapGuard.Lock()
	costTaskMap[t] = t
	costTaskMapGuard.Unlock()
}

func removeCostTask(t *Task) {
	costTaskMapGuard.Lock()
	delete(costTaskMap, t)
	costTaskMapGuard.Unlock()
}

func costTaskLen() int {
	costTaskMapGuard.Lock()
	defer costTaskMapGuard.Unlock()
	return len(costTaskMap)
}

func waitExit() {

	endChan := make(chan bool)

	go func() {

		for {

			l := runningTaskLen()

			if l == 0 {
				break
			}

		}

		endChan <- true

	}()

	<-endChan

}
