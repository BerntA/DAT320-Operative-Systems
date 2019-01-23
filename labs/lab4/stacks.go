// +build !solution

// Leave an empty line above this comment.
package lab4

type UnsafeStack struct {
	top  *Element
	size int
}

func (us *UnsafeStack) Len() int {
	return us.size
}

func (us *UnsafeStack) Push(value interface{}) {
	us.top = &Element{value, us.top}
	us.size++
}

func (us *UnsafeStack) Pop() (value interface{}) {
	if us.size > 0 {
		value, us.top = us.top.value, us.top.next
		us.size--
		return
	}

	return nil
}
