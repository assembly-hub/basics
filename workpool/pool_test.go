package workpool

import (
	"fmt"
	"testing"
	"time"
)

func TestNewWorkPool(t *testing.T) {
	wp := NewWorkPool(100, "test", 0, 50)
	for i := 0; i < 100; i++ {
		wp.SubmitJob(&JobBag{
			JobFunc: func(i ...interface{}) {
				if i[0].(int)%2 == 0 {
					// panic(fmt.Errorf("1111111111111"))
				}
				fmt.Println("------------", i[0])
				// time.Sleep(time.Second * 1)

			},
			Params: []interface{}{i},
		})
	}

	for {
		if wp.IsFinished() {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	wp.ShutDownPool()
	fmt.Println("done")
}
