package gTaskMg

import (
	"sync"
)

type GTask struct {
	GQueue
	sync.WaitGroup
	*GManager
}

func (my *GTask) Init(manager *GManager, name string, qsize int) *GTask {
	my.GManager = manager
	my.GQueue.Init(qsize)
	if my.GManager != nil {
		// 注册默认队列。
		my.GManager.RegistQueue(my.GQueue, name)
		return my
	}
	return nil
}

func (my *GTask) Enter() {
	if my.GManager != nil {
		my.GManager.Enter()
	}
	my.WaitGroup.Add(1)
}

func (my *GTask) Exit() {
	my.WaitGroup.Done()
	if my.GManager != nil {
		my.GManager.Exit()
	}
}

func (my *GTask) Join() {
	my.WaitGroup.Wait()
}
