package orbit

import (
	"fmt"
	"hash/fnv"
	"log"
)

// Worker 接口实现
type Worker interface {
	GetWorkerPoolSize() int
	UseWorkerPool()
	UseSingleWorker(wid int, taskQueue chan *Context)
	JoinTaskQueue(ctx *Context)
}

// worker 结构体
type worker struct {
	poolSize  int
	taskLen   int
	taskQueue []chan *Context
	router    Router
}

// GetWorkerPoolSize 获取工作池大小
func (w *worker) GetWorkerPoolSize() int {
	return w.poolSize
}

// UseWorkerPool 使用 worker pool
func (w *worker) UseWorkerPool() {
	log.Println(fmt.Sprintf("[ WORKER ] worker pool init, size: %d, task length: %d", w.poolSize, w.taskLen))
	for i := 0; i < w.poolSize; i++ {
		w.taskQueue[i] = make(chan *Context, w.taskLen)
		go w.UseSingleWorker(i, w.taskQueue[i])
	}
}

// UseSingleWorker 使用单个 worker
func (w *worker) UseSingleWorker(wid int, taskQueue chan *Context) {
	log.Println(fmt.Sprintf("[ WORKER ] worker %d is ready", wid))
	for {
		select {
		case task := <-taskQueue:
			w.router.do(task)
		}
	}
}

// JoinTaskQueue 加入任务队列
func (w *worker) JoinTaskQueue(ctx *Context) {
	h := fnv.New32a()
	h.Write([]byte(ctx.RemoteAddr()))
	i := int(h.Sum32()) % w.poolSize

	log.Println(fmt.Sprintf("[ WORKER ] worker %d serves for %s, protocol: %d, data: %s", i, ctx.RemoteAddr(), ctx.Protocol(), ctx.RawData()))
	w.taskQueue[i] <- ctx
}
