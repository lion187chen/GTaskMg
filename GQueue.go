package gTaskMg

import (
	"errors"
	"time"
)

type GQueue chan interface{}

func (my *GQueue) Init(size int) *GQueue {
	*my = make(chan interface{}, size)
	return my
}

func (my *GQueue) EnQueueSync(itm interface{}) {
	*my <- itm
}

func (my *GQueue) EnQueue(itm interface{}, timeout time.Duration) error {
	select {
	case *my <- itm:
		return nil
	case <-time.After(timeout):
		return errors.New(ERR_GQUEUE_TIMEOUT)
	}
}

func (my *GQueue) DeQueueSync() interface{} {
	itm := <-*my
	return itm
}

func (my *GQueue) DeQueue(timeout time.Duration) (interface{}, error) {
	var itm interface{}
	select {
	case itm = <-*my:
		return itm, nil
	case <-time.After(timeout):
		return nil, errors.New(ERR_GQUEUE_TIMEOUT)
	}
}

func (my *GQueue) Close() {
	close(*my)
}
