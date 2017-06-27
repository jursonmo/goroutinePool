package goropool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Task struct {
	handle func(...interface{}) interface{}
	args []interface{}
}

type GoRoPool struct {
	wg	*sync.WaitGroup
	jobwg	*sync.WaitGroup
	goroutineNum int
	goroutineFree int32	//空闲goroutine 的个数
	taskQueueSize int
	taskQueue	chan Task

}
func NewTask(f func(...interface{}) interface{}, args ...interface{}) Task {
	return Task {
		handle : f,
		args : args,
	}
}
func doWork(gp *GoRoPool, goroutineId int) {
	gp.wg.Done()
	for task := range gp.taskQueue {
		fmt.Printf("goroutine id:%d get job to do\n", goroutineId)		
		atomic.AddInt32(&gp.goroutineFree, -1)
		task.handle(task.args...)
		atomic.AddInt32(&gp.goroutineFree, 1)
		gp.jobwg.Done()
	}
}

func NewGoRoPool(gn , ts int) *GoRoPool{
	return &GoRoPool {
		goroutineNum : gn,
		goroutineFree : int32(gn),
		taskQueueSize : ts,
	}
}

func (gp *GoRoPool) Run(){
	gp.wg = new(sync.WaitGroup)
	gp.jobwg = new(sync.WaitGroup)
	gp.taskQueue = make( chan Task, gp.taskQueueSize)
	for i := 0; i < gp.goroutineNum; i++ {
		gp.wg.Add(1)
		go doWork(gp, i)
	}
	//wait to create goroutine pool finish
	gp.wg.Wait()
}

func (gp *GoRoPool) AddTask( t Task){
	gp.jobwg.Add(1)
	gp.taskQueue <- t
}

func (gp *GoRoPool) WaitJobDone(){
	gp.jobwg.Wait()
}

func (gp *GoRoPool) GetFree() int{
	return int(atomic.LoadInt32(&gp.goroutineFree))
}
