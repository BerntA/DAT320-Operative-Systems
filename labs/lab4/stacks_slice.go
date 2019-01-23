// +build !solution

// Leave an empty line above this comment.
package lab4

import (
	"sync"
)

const DefaultCap = 10

type SliceStack struct {
	slice []interface{}
	top   int
	lock  sync.Mutex
}

func NewSliceStack() *SliceStack {
	return &SliceStack{
		slice: make([]interface{}, DefaultCap),
		top:   -1,
	}
}

func (ss *SliceStack) Len() int {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	return (ss.top + 1)
}

func (ss *SliceStack) Push(value interface{}) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.top++
	if ss.top == len(ss.slice) {
		// Reallocate
		newSlice := make([]interface{}, len(ss.slice)*2)
		copy(newSlice, ss.slice)
		ss.slice = newSlice
	}

	ss.slice[ss.top] = value
}

func (ss *SliceStack) Pop() (value interface{}) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	var v interface{}
	if ss.top > -1 {
		v = ss.slice[ss.top]
		ss.top--
	}

	return v
}
