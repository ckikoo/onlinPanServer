package timer

import (
	"container/heap"
	"time"
)

type Timer int

func NewTimer() *Timer {
	return new(Timer)
}

// Schedule 安排一个函数在指定时间执行
func (timer *Timer) Push(th *TimeHeap, t time.Time, action func()) {
	heap.Push(th, &TimeItem{priority: t, action: action})
}

func (timer *Timer) Run(th *TimeHeap) {
	for {
		// 检查堆顶元素的时间是否到了
		now := time.Now()
		item := th.Peek()
		if item == nil || !now.After(item.priority) {
			time.Sleep(500 * time.Millisecond) // 没有任务或还未到时间，稍等一会
			continue
		}

		// 时间到了，执行相应的函数
		item = heap.Pop(th).(*TimeItem)
		item.action()
	}
}
