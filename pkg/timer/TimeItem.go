package timer

import (
	"container/heap"
	"time"
)

// TimeItem 代表时间元素
type TimeItem struct {
	Key      interface{} // 键
	Deadline time.Time   // 截止时间
	Action   func()      // 要执行的函数
	Index    int         // 在堆中的索引
}

// TimeHeap 实现了堆接口，用于存储时间元素
type TimeHeap []*TimeItem

// Len 返回堆的长度
func (th TimeHeap) Len() int { return len(th) }

// Less 比较堆中两个元素的时间先后顺序
func (th TimeHeap) Less(i, j int) bool {
	return th[i].Deadline.Before(th[j].Deadline)
}

// Swap 交换堆中两个元素的位置
func (th TimeHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
	th[i].Index = i
	th[j].Index = j
}

// Push 向堆中添加元素
func (th *TimeHeap) Push(x interface{}) {
	item := x.(*TimeItem)
	item.Index = len(*th)
	*th = append(*th, item)
}

// Pop 从堆中移除并返回最小的元素
func (th *TimeHeap) Pop() interface{} {
	old := *th
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	item.Index = -1 // 表示该元素已经不在堆中
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

// PushItem 向堆中添加元素
func (th *TimeHeap) PushItem(key interface{}, deadline time.Time, action func()) {
	item := &TimeItem{
		Key:      key,
		Deadline: deadline,
		Action:   action,
	}
	heap.Push(th, item)
}

// PopItem 从堆中移除并返回最小的元素
func (th *TimeHeap) PopItem() *TimeItem {
	if len(*th) == 0 {
		return nil
	}
	return heap.Pop(th).(*TimeItem)
}

// RemoveItemByKey 根据键从堆中删除指定元素
func (th *TimeHeap) RemoveItemByKey(key interface{}) *TimeItem {
	for i, item := range *th {
		if item.Key == key {
			heap.Remove(th, i)
			return item
		}
	}
	return nil
}
