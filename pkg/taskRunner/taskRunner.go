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

// TaskRunner ...
//
// Whenever we call "WithPriority(...)" a new task will be pushed in the "pushCh" channel and then the code will block
// until the task can actually be executed.
//
// The task runner starts 2 threads.
// - One that receive new tasks from the "WithPriority" function and put them in the priority queue. It then
//   notifies the "popCh" that a new task has been added.
// - The second thread pop tasks from the PQ. Then notify the task that it can be processed. This will unblock the
//   code at "WithPriority", then it waits until that task is done being processed. When the task is completed,
//   it will wait for another one from the "popCh".
//
// This way we can ensure that we ever only have 1 task being executed at the time, but we can queue as many
// as we want with different priorities.
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
	chanLen := 100
	r := &TaskRunner[T]{}
	r.factory = factory
	r.tasks = NewPriorityQueue[*item]()
	r.tasksPushCh = make(chan *item, chanLen)
	r.tasksPopCh = make(chan struct{}, chanLen)
	r.ctx = ctx
	r.start()
	return r
}

func (r *TaskRunner[T]) start() {
	go func() {
		for {
			var t *item
			select {
			case t = <-r.tasksPushCh:
			case <-r.ctx.Done():
				return
			}
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
		for {
			select {
			case <-r.tasksPopCh:
			case <-r.ctx.Done():
				return
			}
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
	task := &item{
		priority:         priority,
		canBeProcessedCh: canBeProcessedCh,
		isDoneCh:         taskIsDoneCh,
	}
	r.tasksPushCh <- task
	select {
	case <-canBeProcessedCh:
	case <-r.ctx.Done():
	}
	t := r.factory()
	t.SetTaskDoneCh(taskIsDoneCh)
	return t
}

// TasksOverview overview of tasks in heap
type TasksOverview struct {
	Low       int64
	Normal    int64
	Important int64
	Critical  int64
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
