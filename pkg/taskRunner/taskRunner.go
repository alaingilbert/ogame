package taskRunner

import (
	"context"
	"sync"
)

type Priority int64

// Priorities
const (
	Low Priority = iota + 1
	Normal
	Important
	Critical
)

// item ...
type item struct {
	canBeProcessedCh chan struct{}
	isDoneCh         chan struct{}
	priority         Priority
	index            int // The index of the item in the heap.
}

func (i *item) GetPriority() int { return int(i.priority) }
func (i *item) GetIndex() int    { return i.index }
func (i *item) SetIndex(idx int) { i.index = idx }

type TaskRunner[T ITask] struct {
	tasks       *PriorityQueue[*item]
	tasksLock   sync.Mutex
	tasksPushCh chan *item
	tasksPopCh  chan struct{}
	factory     func() T
	ctx         context.Context
}

type ITask interface {
	SetTaskDoneCh(ch chan struct{})
}

func NewTaskRunner[T ITask](ctx context.Context, factory func() T) *TaskRunner[T] {
	r := &TaskRunner[T]{}
	r.factory = factory
	r.tasks = NewPriorityQueue[*item]()
	r.tasksPushCh = make(chan *item, 100)
	r.tasksPopCh = make(chan struct{}, 100)
	r.ctx = ctx
	r.start()
	return r
}

func (r *TaskRunner[T]) start() {
	go func() {
		for t := range r.tasksPushCh {
			r.tasksLock.Lock()
			r.tasks.Push(t)
			r.tasksLock.Unlock()
			select {
			case r.tasksPopCh <- struct{}{}:
			case <-r.ctx.Done():
				return
			}
		}
	}()
	go func() {
		for range r.tasksPopCh {
			r.tasksLock.Lock()
			task := r.tasks.Pop()
			r.tasksLock.Unlock()
			close(task.canBeProcessedCh)
			select {
			case <-task.isDoneCh:
			case <-r.ctx.Done():
				return
			}
		}
	}()
}

func (r *TaskRunner[T]) WithPriority(priority Priority) T {
	canBeProcessedCh := make(chan struct{})
	taskIsDoneCh := make(chan struct{})
	task := new(item)
	task.priority = priority
	task.canBeProcessedCh = canBeProcessedCh
	task.isDoneCh = taskIsDoneCh
	r.tasksPushCh <- task
	<-canBeProcessedCh
	t := r.factory()
	t.SetTaskDoneCh(taskIsDoneCh)
	return t
}

// TasksOverview overview of tasks in heap
type TasksOverview struct {
	Low       Priority
	Normal    Priority
	Important Priority
	Critical  Priority
	Total     int64
}

func (r *TaskRunner[T]) GetTasks() (out TasksOverview) {
	r.tasksLock.Lock()
	out.Total = int64(r.tasks.Len())
	for _, item := range r.tasks.Items() {
		switch item.priority {
		case Low:
			out.Low++
		case Normal:
			out.Normal++
		case Important:
			out.Important++
		case Critical:
			out.Critical++
		}
	}
	r.tasksLock.Unlock()
	return
}
