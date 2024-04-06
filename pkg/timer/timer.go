package timer

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// TimerManager 定时器管理器
type TimerManager struct {
	timer       *time.Timer
	timerHeap   *TimeHeap
	closeChan   chan struct{}
	hasItemChan chan struct{}
	mu          *sync.Mutex
}

// NewTimerManager 创建一个新的定时器管理器
func NewTimerManager() (*TimerManager, func()) {
	timerHeap := &TimeHeap{}
	heap.Init(timerHeap)

	timer := &TimerManager{
		timer:       time.NewTimer(0), // 初始化一个立即过期的定时器
		timerHeap:   timerHeap,
		closeChan:   make(chan struct{}),
		hasItemChan: make(chan struct{}, 1), // 使用缓冲通道，避免阻塞
		mu:          &sync.Mutex{},
	}

	go timer.Run()
	return timer, timer.Close
}

// Run 启动定时器管理器
func (tm *TimerManager) Run() {
	defer func() {
		if tm.timer != nil {
			tm.timer.Stop()
		}
		close(tm.hasItemChan)
	}()

	for {
		select {
		case <-tm.closeChan:
			return

		case <-tm.timer.C:
			tm.executeExpiredTasks()
		case <-tm.hasItemChan:
			tm.resetTimer()
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// resetTimer 重新设置定时器
func (tm *TimerManager) resetTimer() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if tm.timerHeap.Len() > 0 {
		duration := tm.timerHeap.Peek().Deadline.Sub(time.Now())
		tm.timer.Reset(duration)
	}
}

// executeExpiredTasks 执行到期任务
func (tm *TimerManager) executeExpiredTasks() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for tm.timerHeap.Len() > 0 {
		nextItem := tm.timerHeap.Peek()
		if nextItem.Deadline.After(time.Now()) {
			return
		}
		heap.Pop(tm.timerHeap).(*TimeItem).Action()
	}
}

// Add 添加一个定时任务
func (tm *TimerManager) Add(key interface{}, deadline time.Time, action func()) {

	item := &TimeItem{
		Key:      key,
		Deadline: deadline,
		Action:   action,
	}
	tm.mu.Lock()

	heap.Push(tm.timerHeap, item)
	tm.mu.Unlock()

	tm.hasItemChan <- struct{}{}
}

// Del 根据键删除一个定时任务
func (tm *TimerManager) Del(key interface{}) *TimeItem {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for i, item := range *tm.timerHeap {
		if item.Key == key {
			fmt.Println("1111")
			return heap.Remove(tm.timerHeap, i).(*TimeItem)
		}
	}
	return nil
}

// Close 关闭定时器管理器
func (tm *TimerManager) Close() {
	close(tm.closeChan)
}
