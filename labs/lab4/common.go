package lab4

type Stack interface {
	Len() int
	Push(value interface{})
	Pop() interface{}
}

type Element struct {
	value interface{}
	next  *Element
}
