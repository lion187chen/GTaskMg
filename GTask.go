package gTaskMg

import (
	"sync"
)

type GTask struct {
	runner interface{}
	name   string
	GQueue
	sync.WaitGroup
	*GManager
}

func (obj *GTask) Init(runner interface{}, name string, manager *GManager, qsize int) *GTask {
	obj.runner = runner
	obj.name = name
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

func (obj *GTask) CallRunner() interface{} {
	obj.Enter()
	return obj.runner
}

func (obj *GTask) Runner() interface{} {
	return obj.runner
}

func (obj *GTask) Name() string {
	return obj.name
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
