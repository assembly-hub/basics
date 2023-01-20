// Package workpool 定义协程池
package workpool

type WorkFunc func(...interface{})

type JobBag struct {
	JobFunc WorkFunc
	Params  []interface{}
}
