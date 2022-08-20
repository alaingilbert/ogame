package taskRunner

import (
	"fmt"
	"time"
)

type testItem struct {
	taskDoneCh chan struct{}
}

func (i *testItem) SetTaskDoneCh(ch chan struct{}) {
	i.taskDoneCh = ch
}

func (i *testItem) DoSomething(name string) {
	defer close(i.taskDoneCh)
	time.Sleep(250 * time.Millisecond)
	fmt.Println("Do something " + name)
}

func (i *testItem) DoSomethingElse(name string) {
	defer close(i.taskDoneCh)
	time.Sleep(250 * time.Millisecond)
	fmt.Println("Do something else " + name)
}

//func TestA(t *testing.T) {
//	wg := &sync.WaitGroup{}
//	wg.Add(6)
//	factory := func() *testItem { return &testItem{} }
//	tr := NewTaskRunner[*testItem](context.Background(), factory)
//	tr.WithPriority(Low).DoSomething("Always first")
//	go func() { tr.WithPriority(Low).DoSomething("A"); wg.Done() }()
//	go func() { tr.WithPriority(Low).DoSomethingElse("B"); wg.Done() }()
//	go func() { tr.WithPriority(Low).DoSomething("C"); wg.Done() }()
//	go func() { tr.WithPriority(Low).DoSomething("D"); wg.Done() }()
//	go func() { time.Sleep(480 * time.Millisecond); tr.WithPriority(Critical).DoSomething("E"); wg.Done() }()
//	go func() { time.Sleep(470 * time.Millisecond); tr.WithPriority(Important).DoSomething("F"); wg.Done() }()
//	wg.Wait()
//}
