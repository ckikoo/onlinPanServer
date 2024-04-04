package timer

import (
	"time"
)

type TimeItem struct {
	priority time.Time // 优先级，即时间
	action   func()    // 要执行的函数
}

// TimeHeap 实现了堆接口，用于存储时间元素
type TimeHeap []*TimeItem

// Len 返回堆的长度
func (th TimeHeap) Len() int { return len(th) }

// Less 比较堆中两个元素的时间先后顺序
func (th TimeHeap) Less(i, j int) bool {
	return th[i].priority.Before(th[j].priority)
}

// Swap 交换堆中两个元素的位置
func (th TimeHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
}

// Push 向堆中添加元素
func (th *TimeHeap) Push(x interface{}) {
	item := x.(*TimeItem)
	*th = append(*th, item)
}

// Pop 从堆中移除并返回最小的元素
func (th *TimeHeap) Pop() interface{} {
	old := *th
	n := len(old)
	item := old[n-1]
	*th = old[0 : n-1]
	return item
}

// Peek 返回堆顶元素
func (th TimeHeap) Peek() *TimeItem {
	if len(th) == 0 {
		return nil
	}
	return th[0]
}
