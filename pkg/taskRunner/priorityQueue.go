package taskRunner

type IPQItem interface {
	GetIndex() int
	SetIndex(idx int)
	GetPriority() int
}

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
