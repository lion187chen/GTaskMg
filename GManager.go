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

func (obj *GManager) Init() *GManager {
	obj.tasks = make(map[string]*GTask)
	obj.queues = make(map[string]GQueue)
	return obj
}

func (obj *GManager) RegistTask(task *GTask, name string) {
	obj.tlock.Lock()
	defer obj.tlock.Unlock()
	obj.tasks[name] = task
}

func (obj *GManager) DeleteTask(name string) {
	obj.tlock.RLock()
	defer obj.tlock.RUnlock()

	_, f := obj.tasks[name]
	if f {
		obj.DeleteQueue(name)
		delete(obj.tasks, name)
	}
}

func (obj *GManager) CreateTask(name string, qsize int) *GTask {
	t := new(GTask).Init(obj, name, qsize)
	return t
}

func (obj *GManager) RegistQueue(queue GQueue, name string) {
	obj.qlock.RLock()
	defer obj.qlock.RUnlock()
	obj.queues[name] = queue
}

func (obj *GManager) DeleteQueue(name string) {
	obj.qlock.RLock()
	defer obj.qlock.RUnlock()

	q, f := obj.queues[name]
	if f {
		q.Close()
		delete(obj.queues, name)
	}
}

func (obj *GManager) CreateQueue(name string, qsize int) *GQueue {
	q := new(GQueue)
	q.Init(qsize)
	obj.RegistQueue(*q, name)
	return q
}

func (obj *GManager) Enter() {
	obj.WaitGroup.Add(1)
}

func (obj *GManager) Exit() {
	obj.WaitGroup.Done()
}

func (obj *GManager) Join() {
	obj.WaitGroup.Wait()
}

func (obj *GManager) Broadcast(event interface{}) {
	obj.tlock.RLock()
	defer obj.tlock.RUnlock()
	// Last to first.
	for name := range obj.tasks {
		t := obj.tasks[name]
		t.GQueue.EnQueueSync(event)
	}
}

func (obj *GManager) BroadcastWithout(event interface{}, without string) {
	obj.tlock.RLock()
	defer obj.tlock.RUnlock()
	// Last to first.
	for name := range obj.tasks {
		if without == name {
			continue
		}
		t := obj.tasks[name]
		t.GQueue.EnQueueSync(event)
	}
}

func (obj *GManager) ReqExit() {
	obj.Broadcast(GMSG_EXIT)
}

func (obj *GManager) ReqTaskExit(name string) {
	t, f := obj.tasks[name]
	if f {
		t.GQueue.EnQueueSync(GMSG_EXIT)
	}
}

// 这是一个比较危险的函数，在使用时需要确保对应的任务不会被释放。
// 如果只是向队列或者任务发生消息，应使用：GManager.EnQueueSync() 或者 GManager.EnQueue()
func (obj *GManager) GetGTask(name string) (*GTask, bool) {
	obj.tlock.RLock()
	defer obj.tlock.RUnlock()

	t, f := obj.tasks[name]
	return t, f
}

// 这是一个比较危险的函数，在使用时需要确保对应的队列不会被释放，或对通道是否关闭进行检查。
// 如果只是向队列或者任务发生消息，应使用：GManager.EnQueueSync() 或者 GManager.EnQueue()
func (obj *GManager) GetQueue(name string) (GQueue, bool) {
	obj.qlock.RLock()
	defer obj.qlock.RUnlock()

	q, f := obj.queues[name]
	return q, f
}

func (obj *GManager) EnQueueSync(name string, itm interface{}) error {
	obj.qlock.RLock()
	defer obj.qlock.RUnlock()

	q, f := obj.queues[name]
	if !f {
		return errors.New(ERR_GQUEUE_NOFOUND)
	}
	q.EnQueueSync(itm)
	return nil
}

func (obj *GManager) EnQueue(name string, itm interface{}, timeout time.Duration) error {
	obj.qlock.RLock()
	defer obj.qlock.RUnlock()

	q, f := obj.queues[name]
	if !f {
		return errors.New(ERR_GQUEUE_NOFOUND)
	}
	return q.EnQueue(itm, timeout)
}
