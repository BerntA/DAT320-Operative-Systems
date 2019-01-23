// +build !solution

// Leave an empty line above this comment.
package lab5

type FIFOQueue struct {
	q     []interface{}
	size  int
	head  int
	tail  int
	count int
}

func NewFIFOQueue(size int) *FIFOQueue {
	return &FIFOQueue{q: make([]interface{}, size), size: size}
}

func (fq *FIFOQueue) Enqueue(value interface{}) {
	if fq.head == fq.tail && fq.count > 0 {
		q := make([]interface{}, len(fq.q)+fq.size)
		copy(q, fq.q[fq.head:])
		copy(q[len(fq.q)-fq.head:], fq.q[:fq.head])
		fq.head = 0
		fq.tail = len(fq.q)
		fq.q = q
	}
	fq.q[fq.tail] = value
	fq.tail = (fq.tail + 1) % len(fq.q)
	fq.count++
}

func (fq *FIFOQueue) Dequeue() interface{} {
	if fq.count == 0 {
		return nil
	}

	value := fq.q[fq.head]
	fq.head = (fq.head + 1) % len(fq.q)
	fq.count--
	return value
}

func (fq *FIFOQueue) Flush() {
	fq.q = make([]interface{}, fq.size)
	fq.head = 0
	fq.tail = 0
	fq.count = 0
}

func (fq *FIFOQueue) Empty() bool {
	return fq.count == 0
}

func (fq *FIFOQueue) Len() int {
	return fq.count
}
