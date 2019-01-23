// +build !solution

// Leave an empty line above this comment.
package lab5

import (
	"sync"
)

type CondQueue struct {
	queue     *FIFOQueue
	condEmpty *sync.Cond
	lock      sync.Mutex
}

func NewCondQueue(size int) *CondQueue {
	condQueue := CondQueue{lock: sync.Mutex{}}
	condQueue.condEmpty = sync.NewCond(&condQueue.lock)
	condQueue.queue = NewFIFOQueue(size)
	return &condQueue
}

func (cq *CondQueue) Enqueue(value interface{}) {
	cq.lock.Lock()
	defer cq.lock.Unlock()

	cq.queue.Enqueue(value)
	cq.condEmpty.Signal()
}

func (cq *CondQueue) Dequeue() interface{} {
	cq.lock.Lock()
	defer cq.lock.Unlock()

	for cq.queue.Empty() {
		cq.condEmpty.Wait()
	}

	return cq.queue.Dequeue()
}

func (cq *CondQueue) Flush() {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	cq.queue.Flush()
}

func (cq *CondQueue) Empty() bool {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	return cq.queue.Empty()
}

func (cq *CondQueue) Len() int {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	return cq.queue.Len()
}
