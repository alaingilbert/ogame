package mtx

import (
	"cmp"
	"sync"
)

type Mtx[T any] struct {
	sync.Mutex
	v T
}

func NewMtx[T any](v T) Mtx[T] {
	return Mtx[T]{v: v}
}

func (m *Mtx[T]) Val() *T {
	return &m.v
}

func (m *Mtx[T]) Get() T {
	m.Lock()
	defer m.Unlock()
	return m.v
}

func (m *Mtx[T]) Set(v T) {
	m.Lock()
	defer m.Unlock()
	m.v = v
}

func (m *Mtx[T]) With(clb func(v *T)) {
	_ = m.WithE(func(tx *T) error {
		clb(tx)
		return nil
	})
}

func (m *Mtx[T]) WithE(clb func(v *T) error) error {
	m.Lock()
	defer m.Unlock()
	return clb(&m.v)
}

//----------------------

type RWMtx[T any] struct {
	sync.RWMutex
	v T
}

func New[T any](v T) RWMtx[T] {
	return RWMtx[T]{v: v}
}

func (m *RWMtx[T]) Val() *T {
	return &m.v
}

func (m *RWMtx[T]) Get() T {
	m.RLock()
	defer m.RUnlock()
	return m.v
}

func (m *RWMtx[T]) Set(v T) {
	m.Lock()
	defer m.Unlock()
	m.v = v
}

func (m *RWMtx[T]) Replace(newVal T) (old T) {
	m.With(func(v *T) {
		old = *v
		*v = newVal
	})
	return
}

func (m *RWMtx[T]) RWith(clb func(v T)) {
	_ = m.RWithE(func(tx T) error {
		clb(tx)
		return nil
	})
}

func (m *RWMtx[T]) RWithE(clb func(v T) error) error {
	m.RLock()
	defer m.RUnlock()
	return clb(m.v)
}

func (m *RWMtx[T]) With(clb func(v *T)) {
	_ = m.WithE(func(tx *T) error {
		clb(tx)
		return nil
	})
}

func (m *RWMtx[T]) WithE(clb func(v *T) error) error {
	m.Lock()
	defer m.Unlock()
	return clb(&m.v)
}

//----------------------

type RWMtxMap[K cmp.Ordered, V any] struct {
	RWMtx[map[K]V]
}

func NewMap[K cmp.Ordered, V any]() RWMtxMap[K, V] {
	return RWMtxMap[K, V]{RWMtx: New(make(map[K]V))}
}

func (m *RWMtxMap[K, V]) SetKey(k K, v V) {
	m.With(func(m *map[K]V) { (*m)[k] = v })
}

func (m *RWMtxMap[K, V]) GetKey(k K) (out V, ok bool) {
	m.RWith(func(m map[K]V) { out, ok = m[k] })
	return
}

func (m *RWMtxMap[K, V]) HasKey(k K) (found bool) {
	m.RWith(func(m map[K]V) { _, found = m[k] })
	return
}

func (m *RWMtxMap[K, V]) TakeKey(k K) (out V, ok bool) {
	m.With(func(m *map[K]V) {
		out, ok = (*m)[k]
		if ok {
			delete(*m, k)
		}
	})
	return
}

func (m *RWMtxMap[K, V]) DeleteKey(k K) {
	m.With(func(m *map[K]V) { delete(*m, k) })
	return
}

func (m *RWMtxMap[K, V]) Each(clb func(K, V)) {
	m.RWith(func(m map[K]V) {
		for k, v := range m {
			clb(k, v)
		}
	})
}

//----------------------

type RWMtxSlice[T any] struct {
	RWMtx[[]T]
}

// Clear clears the slice, removing all values
func (s *RWMtxSlice[T]) Clear() {
	s.With(func(v *[]T) { *v = nil; *v = make([]T, 0) })
}

// Remove removes the element at position i within the slice,
// shifting all elements after it to the left
// Panics if index is out of bounds
func (s *RWMtxSlice[T]) Remove(i int) (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[i], (*v)[:i+copy((*v)[i:], (*v)[i+1:])] })
	return
}

func (s *RWMtxSlice[T]) Each(clb func(T)) {
	s.RWith(func(v []T) {
		for _, e := range v {
			clb(e)
		}
	})
}

func (s *RWMtxSlice[T]) Append(els ...T) {
	s.With(func(v *[]T) { *v = append(*v, els...) })
}

func (s *RWMtxSlice[T]) Unshift(el T) {
	s.With(func(v *[]T) { *v = append([]T{el}, *v...) })
}

func (s *RWMtxSlice[T]) Clone() (out []T) {
	s.RWith(func(v []T) {
		out = make([]T, len(v))
		copy(out, v)
	})
	return
}

// Len returns the length of the slice
func (s *RWMtxSlice[T]) Len() (out int) {
	s.RWith(func(v []T) { out = len(v) })
	return
}

//----------------------

type RWMtxUInt64[T ~uint64] struct {
	RWMtx[T]
}

func (s *RWMtxUInt64[T]) Incr(diff T) {
	s.With(func(v *T) { *v += diff })
}
