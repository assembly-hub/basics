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
	quit              chan bool
	executeIntervalMS int64
	safeFunc          func(job *JobBag)
}

func newWork(workID string, executeIntervalMS int64) *worker {
	w := new(worker)
	w.WorkID = workID
	w.quit = make(chan bool)
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
		for {
			wp.WorkerQueue <- w
			select {
			case job := <-w.jobData:
				// logtool.Printf("worker: %s, will execute taskfunc.\n", w.WorkID)
				w.safeFunc(job)
			case q := <-w.quit:
				if q {
					// logtool.Printf("worker: %s, will stop.\n", w.WorkID)
					return
				}
			}
			if w.executeIntervalMS > 0 {
				time.Sleep(time.Duration(w.executeIntervalMS) * time.Millisecond)
			}
		}
	}()
}

func (w *worker) stopWorker() {
	w.quit <- true
	close(w.jobData)
}
