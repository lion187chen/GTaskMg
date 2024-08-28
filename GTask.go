package gTaskMg

import (
	"sync"
)

type GTask struct {
	GQueue
	sync.WaitGroup
	*GManager
}

func (obj *GTask) Init(manager *GManager, name string, qsize int) *GTask {
	obj.GManager = manager
	obj.GQueue.Init(qsize)
	if obj.GManager != nil {
		// 注册 Task。
		obj.GManager.RegistTask(obj, name)
		// 注册默认队列。
		obj.GManager.RegistQueue(obj.GQueue, name)
		return obj
	}
	return nil
}

func (obj *GTask) Enter() {
	if obj.GManager != nil {
		obj.GManager.Enter()
	}
	obj.WaitGroup.Add(1)
}

func (obj *GTask) Exit() {
	obj.WaitGroup.Done()
	if obj.GManager != nil {
		obj.GManager.Exit()
	}
}

func (obj *GTask) Join() {
	obj.WaitGroup.Wait()
}
