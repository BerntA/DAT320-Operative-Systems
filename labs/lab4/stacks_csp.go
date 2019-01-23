// +build !solution

// Leave an empty line above this comment.
package lab4

// Request Types
const (
	REQUEST_LEN  = 1
	REQUEST_PUSH = 2
	REQUEST_POP  = 3
)

type Response struct {
	v interface{}
}

type Request struct {
	rtype    int
	v        interface{}
	response chan Response
}

type CspStack struct {
	top     *Element
	size    int
	request chan Request
}

func NewCspStack() *CspStack {
	var cspStack = &CspStack{request: make(chan Request)}
	go cspStack.run()
	return cspStack
}

func (cs *CspStack) Len() int {
	res := make(chan Response)
	cs.request <- Request{rtype: REQUEST_LEN, response: res}
	size, ok := (<-res).v.(int)
	if ok {
		return size
	}
	return -1
}

func (cs *CspStack) Push(value interface{}) {
	res := make(chan Response)
	cs.request <- Request{rtype: REQUEST_PUSH, v: value, response: res}
	<-res // Wait for response.
}

func (cs *CspStack) Pop() (value interface{}) {
	res := make(chan Response)
	cs.request <- Request{rtype: REQUEST_POP, response: res}
	return (<-res).v
}

func (cs *CspStack) run() {
	for {
		req := <-cs.request
		if req.response == nil {
			continue
		}

		switch req.rtype {
		case REQUEST_LEN:
			req.response <- Response{v: cs.size}

		case REQUEST_PUSH:
			cs.top = &Element{req.v, cs.top}
			cs.size++
			req.response <- Response{}

		case REQUEST_POP:
			if cs.size > 0 {
				oldTop := cs.top.value
				cs.top = cs.top.next
				cs.size--
				req.response <- Response{v: oldTop}
			} else {
				req.response <- Response{}
			}
		}
	}
}
