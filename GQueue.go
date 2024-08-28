package gTaskMg

import (
	"errors"
	"time"
)

type GQueue chan interface{}

func (obj *GQueue) Init(size int) *GQueue {
	*obj = make(chan interface{}, size)
	return obj
}

func (obj *GQueue) EnQueueSync(itm interface{}) {
	*obj <- itm
}

func (obj *GQueue) EnQueue(itm interface{}, timeout time.Duration) error {
	select {
	case *obj <- itm:
		return nil
	case <-time.After(timeout):
		return errors.New(ERR_GQUEUE_TIMEOUT)
	}
}

func (obj *GQueue) DeQueueSync() interface{} {
	itm := <-*obj
	return itm
}

func (obj *GQueue) DeQueue(timeout time.Duration) (interface{}, error) {
	var itm interface{}
	select {
	case itm = <-*obj:
		return itm, nil
	case <-time.After(timeout):
		return nil, errors.New(ERR_GQUEUE_TIMEOUT)
	}
}

func (obj *GQueue) Close() {
	close(*obj)
}
