// Package workpool 定义协程池
package workpool

import (
	"fmt"
)

type data struct {
	WorkerQueue     chan *worker
	WorkerList      []*worker
	maxPoolSize     int
	poolName        string
	jobQueue        chan *JobBag
	jobQueueMaxSize int
	quit            chan struct{}
	isShutDown      bool
	finishNotify    chan struct{}
}

type WorkPool interface {
	// SubmitJob 提交任务
	SubmitJob(job ...*JobBag)
	// GetPoolIdleSize get pool idle size
	GetPoolIdleSize() int
	// GetJobQueueIdleSize 获取job队列空闲成都
	GetJobQueueIdleSize() int
	// ShutDownPool 关闭工作池
	ShutDownPool()
	// IsShutDownPool 获取是否关闭
	IsShutDownPool() bool
	// IsFinished get pool execute info, true is finished false is running
	IsFinished() bool

	// WaitFinish 等待任务完成通知
	WaitFinish()
}

// NewWorkPool 初始化work pool
// maxPoolSize work pool size > 0
// poolName pool name
// executeIntervalMS worker executeIntervalMS, 0 is not be used
// jobQueueMaxSize job queue max size
func NewWorkPool(maxPoolSize int, poolName string, executeIntervalMS int64, jobQueueMaxSize int) WorkPool {
	if maxPoolSize <= 0 || executeIntervalMS < 0 {
		panic("maxPoolSize must gt 0, and executeIntervalMS gte 0")
	}

	wp := new(data)
	wp.maxPoolSize = maxPoolSize
	wp.poolName = poolName
	wp.isShutDown = false

	wp.WorkerQueue = make(chan *worker, maxPoolSize)
	wp.WorkerList = make([]*worker, 0, maxPoolSize)

	if jobQueueMaxSize >= maxPoolSize {
		wp.jobQueueMaxSize = jobQueueMaxSize
	} else {
		wp.jobQueueMaxSize = maxPoolSize
	}
	wp.jobQueue = make(chan *JobBag, wp.jobQueueMaxSize)

	wp.quit = make(chan struct{})
	wp.finishNotify = make(chan struct{})

	for i := 0; i < maxPoolSize; i++ {
		w := newWork(fmt.Sprintf("%s-%d", poolName, i), executeIntervalMS)
		wp.WorkerList = append(wp.WorkerList, w)
		w.startWorker(wp)
	}

	wp.jobQueueManager()
	return wp
}

func (w *data) WaitFinish() {
	if w.IsFinished() {
		return
	}

	select {
	case <-w.finishNotify:
		return
	}
}

func (w *data) sendFinishNotify() {
	if w.IsFinished() {
		defer func() {
			_ = recover()
		}()
		select {
		case w.finishNotify <- struct{}{}:
		default:
		}
	}
}

func (w *data) jobQueueManager() {
	go func() {
		defer func() {
			_ = recover()
		}()
		for {
			if w.isShutDown {
				return
			}

			select {
			case job, ok := <-w.jobQueue:
				if !ok {
					return
				}
				wk := <-w.WorkerQueue
				wk.jobData <- job
			case <-w.quit:
				return
			}
		}
	}()
}

// GetPoolIdleSize get pool idle size
func (w *data) GetPoolIdleSize() int {
	return len(w.WorkerQueue)
}

// GetJobQueueIdleSize 获取job队列空闲长度
func (w *data) GetJobQueueIdleSize() int {
	return cap(w.jobQueue) - len(w.jobQueue)
}

func (w *data) IsShutDownPool() bool {
	return w.isShutDown
}

// IsFinished get pool execute info, true is finished false is running
func (w *data) IsFinished() bool {
	return len(w.jobQueue) == 0 && len(w.WorkerQueue) == cap(w.WorkerQueue)
}

// SubmitJob 提交任务
func (w *data) SubmitJob(job ...*JobBag) {
	if w.isShutDown {
		panic("already shutdown")
	}

	for _, v := range job {
		w.jobQueue <- v
	}
}

// ShutDownPool 关闭工作池
func (w *data) ShutDownPool() {
	w.isShutDown = true

	w.quit <- struct{}{}
	for i := 0; i < len(w.WorkerList); i++ {
		w.WorkerList[i].stopWorker()
	}

	close(w.quit)
	close(w.WorkerQueue)
	close(w.jobQueue)
	close(w.finishNotify)
}
