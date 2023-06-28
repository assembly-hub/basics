// Package workpool
package workpool

import (
	"log"
	"runtime/debug"
	"time"
)

type worker struct {
	WorkID            string
	jobData           chan *JobBag
	quit              chan struct{}
	executeIntervalMS int64
	safeFunc          func(job *JobBag)
}

func newWork(workID string, executeIntervalMS int64) *worker {
	w := new(worker)
	w.WorkID = workID
	w.quit = make(chan struct{})
	w.jobData = make(chan *JobBag)
	w.executeIntervalMS = executeIntervalMS
	w.safeFunc = func(job *JobBag) {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("worker: %s, error is %v.\n%s", w.WorkID, p, string(debug.Stack()))
			}
		}()

		if job.JobFunc != nil {
			job.JobFunc(job.Params...)
		} else {
			log.Printf("worker: %s, execute taskfunc found some error, msg is taskfunc is nil.\n", w.WorkID)
		}
	}
	return w
}

func (w *worker) startWorker(wp *data) {
	go func() {
		defer func() {
			_ = recover()
		}()
		for {
			wp.WorkerQueue <- w
			wp.sendFinishNotify()
			select {
			case job, ok := <-w.jobData:
				if !ok {
					return
				}
				w.safeFunc(job)
			case <-w.quit:
				return
			}
			if w.executeIntervalMS > 0 {
				time.Sleep(time.Duration(w.executeIntervalMS) * time.Millisecond)
			}
		}
	}()
}

func (w *worker) stopWorker() {
	w.quit <- struct{}{}
	close(w.jobData)
}
