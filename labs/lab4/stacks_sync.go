// +build !solution

// Leave an empty line above this comment.
package lab4

import (
	"sync"
)

type SafeStack struct {
	top  *Element
	size int
	lock sync.Mutex
}

func (ss *SafeStack) Len() int {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	return ss.size
}

func (ss *SafeStack) Push(value interface{}) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	ss.top = &Element{value, ss.top}
	ss.size++
}

func (ss *SafeStack) Pop() (value interface{}) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	if ss.size > 0 {
		value, ss.top = ss.top.value, ss.top.next
		ss.size--
		return
	}

	return nil
}
