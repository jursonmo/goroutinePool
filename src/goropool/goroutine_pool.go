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
	goroutineNum int32
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
func doWork(gp *GoRoPool, goroutineId int32, t Task) {
	
	if t.handle != nil {
		fmt.Printf("goroutineId %d can't wait working....\n", goroutineId)
		t.handle(t.args...)
		gp.jobwg.Done()
	}

	gp.wg.Done()//create goroutine ok and wait to do job
	atomic.AddInt32(&gp.goroutineFree, 1)
	fmt.Printf("goroutineId %d standby....\n", goroutineId)
	for task := range gp.taskQueue {
		fmt.Printf("goroutine id:%d get job to do\n", goroutineId)				
		if task.handle != nil {
			atomic.AddInt32(&gp.goroutineFree, -1)
			fmt.Printf("goroutineId %d working....\n", goroutineId)
			task.handle(task.args...)
			atomic.AddInt32(&gp.goroutineFree, 1)
		}else {
			fmt.Printf("task.handle is nil \n")
		}		
		gp.jobwg.Done()
	}
}

func NewGoRoPool(gn , ts int) *GoRoPool{
	return &GoRoPool {
		goroutineNum : int32(gn),
		goroutineFree : int32(0),
		taskQueueSize : ts,
	}
}

func (gp *GoRoPool) Run(){
	gp.wg = new(sync.WaitGroup)
	gp.jobwg = new(sync.WaitGroup)
	gp.taskQueue = make( chan Task, gp.taskQueueSize)
	for i := int32(1); i <= gp.goroutineNum; i++ {
		gp.wg.Add(1)
		go doWork(gp, i, Task{nil, nil})
	}
	//wait to create goroutine pool finish
	gp.wg.Wait()
}
func (gp *GoRoPool) createOne(t Task){
	 //goroutineId :=atomic.LoadInt32(&gp.goroutineNum)
	goroutineId := atomic.AddInt32(&gp.goroutineNum, 1)
	gp.wg.Add(1)//i don't have to do this
	go doWork(gp, goroutineId, t)
}

func (gp *GoRoPool) AddTask( t Task){
	gp.jobwg.Add(1)
	if atomic.LoadInt32(&gp.goroutineFree) == 0 {
		//there is no free goroutine, so create another one to handle this task		
		gp.createOne(t)
	}else {
		gp.taskQueue <- t
	}
	
}

func (gp *GoRoPool) WaitJobDone(){
	gp.jobwg.Wait()
}

func (gp *GoRoPool) GetFree() int{
	return int(atomic.LoadInt32(&gp.goroutineFree))
}
