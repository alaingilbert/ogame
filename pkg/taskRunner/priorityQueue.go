package taskRunner

import "container/heap"

type IPQItem interface {
	GetIndex() int
	SetIndex(idx int)
	GetPriority() int
}

type PriorityQueue[T IPQItem] struct {
	items priorityQueue[T]
}

func NewPriorityQueue[T IPQItem]() *PriorityQueue[T] {
	pq := new(PriorityQueue[T])
	pq.items = make(priorityQueue[T], 0)
	heap.Init(&pq.items)
	return pq
}

func (p *PriorityQueue[T]) Push(item T) { heap.Push(&p.items, item) }
func (p *PriorityQueue[T]) Pop() T      { return heap.Pop(&p.items).(T) }
func (p *PriorityQueue[T]) Len() int    { return p.items.Len() }
func (p *PriorityQueue[T]) Items() []T  { return p.items }

// A priorityQueue implements heap.Interface and holds Items.
type priorityQueue[T IPQItem] []T

func (pq *priorityQueue[T]) Len() int { return len(*pq) }

func (pq *priorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return (*pq)[i].GetPriority() > (*pq)[j].GetPriority()
}

func (pq *priorityQueue[T]) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].SetIndex(i)
	(*pq)[j].SetIndex(j)
}

// Push ...
func (pq *priorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(T)
	item.SetIndex(n)
	*pq = append(*pq, item)
}

// Pop ...
func (pq *priorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.SetIndex(-1) // for safety
	*pq = old[0 : n-1]
	return item
}
