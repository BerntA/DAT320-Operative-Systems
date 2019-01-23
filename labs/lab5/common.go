package lab5

type Queue interface {
	Enqueue(value interface{})
	Dequeue() interface{}
	Flush()
	Empty() bool
	Len() int
}
