package gTaskMg

import (
	"errors"
	"sync"
	"time"
)

type GManager struct {
	sync.WaitGroup
	tlock  sync.RWMutex
	qlock  sync.RWMutex
	tasks  map[string]*GTask
	queues map[string]GQueue
}

const (
	GMSG_EXIT string = "GManager.Exit"
)

const (
	ERR_GQUEUE_NOFOUND string = "cannot find queue"
	ERR_GQUEUE_TIMEOUT string = "gqueue timeout"
)

func (my *GManager) Init() *GManager {
	my.tasks = make(map[string]*GTask)
	return my
}

func (my *GManager) RegistTask(task *GTask, name string) {
	my.tlock.Lock()
	defer my.tlock.Unlock()
	my.tasks[name] = task
}

func (my *GManager) DeleteTask(name string) {
	my.tlock.RLock()
	defer my.tlock.RUnlock()

	_, f := my.tasks[name]
	if f {
		my.DeleteQueue(name)
		delete(my.tasks, name)
	}
}

func (my *GManager) CreateTask(name string, qsize int) *GTask {
	t := new(GTask).Init(my, name, qsize)
	t.GManager = my
	my.RegistTask(t, name)

	return t
}

func (my *GManager) RegistQueue(queue GQueue, name string) {
	my.qlock.RLock()
	defer my.qlock.RUnlock()
	my.queues[name] = queue
}

func (my *GManager) DeleteQueue(name string) {
	my.qlock.RLock()
	defer my.qlock.RUnlock()

	q, f := my.queues[name]
	if f {
		q.Close()
		delete(my.queues, name)
	}
}

func (my *GManager) CreateQueue(name string, qsize int) *GQueue {
	q := new(GQueue)
	q.Init(qsize)
	my.RegistQueue(*q, name)
	return q
}

func (my *GManager) Enter() {
	my.WaitGroup.Add(1)
}

func (my *GManager) Exit() {
	my.WaitGroup.Done()
}

func (my *GManager) Join() {
	my.WaitGroup.Wait()
}

func (my *GManager) Broadcast(event interface{}) {
	my.tlock.RLock()
	defer my.tlock.RUnlock()
	// Last to first.
	for name := range my.tasks {
		t := my.tasks[name]
		t.GQueue.EnQueueSync(event)
	}
}

func (my *GManager) BroadcastWithout(event interface{}, without string) {
	my.tlock.RLock()
	defer my.tlock.RUnlock()
	// Last to first.
	for name := range my.tasks {
		if without == name {
			continue
		}
		t := my.tasks[name]
		t.GQueue.EnQueueSync(event)
	}
}

func (my *GManager) ReqExit() {
	my.Broadcast(GMSG_EXIT)
}

func (my *GManager) ReqTaskExit(name string) {
	t, f := my.tasks[name]
	if f {
		t.GQueue.EnQueueSync(GMSG_EXIT)
	}
}

// 这是一个比较危险的函数，在使用时需要确保对应的任务不会被释放。
// 如果只是向队列或者任务发生消息，应使用：GManager.EnQueueSync() 或者 GManager.EnQueue()
func (my *GManager) GetGTask(name string) (*GTask, bool) {
	my.tlock.RLock()
	defer my.tlock.RUnlock()

	t, f := my.tasks[name]
	return t, f
}

// 这是一个比较危险的函数，在使用时需要确保对应的队列不会被释放，或对通道是否关闭进行检查。
// 如果只是向队列或者任务发生消息，应使用：GManager.EnQueueSync() 或者 GManager.EnQueue()
func (my *GManager) GetQueue(name string) (GQueue, bool) {
	my.qlock.RLock()
	defer my.qlock.RUnlock()

	q, f := my.queues[name]
	return q, f
}

func (my *GManager) EnQueueSync(name string, itm interface{}) error {
	my.qlock.RLock()
	defer my.qlock.RUnlock()

	q, f := my.queues[name]
	if !f {
		return errors.New(ERR_GQUEUE_NOFOUND)
	}
	q.EnQueueSync(itm)
	return nil
}

func (my *GManager) EnQueue(name string, itm interface{}, timeout time.Duration) error {
	my.qlock.RLock()
	defer my.qlock.RUnlock()

	q, f := my.queues[name]
	if !f {
		return errors.New(ERR_GQUEUE_NOFOUND)
	}
	return q.EnQueue(itm, timeout)
}
