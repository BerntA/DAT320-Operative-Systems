package lab5

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestCondQueue(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	condQueue := NewCondQueue(5000)
	fmt.Println("Test Conditional Queue")
	testConcurrentQueueAccess(condQueue)
}

func TestOpsUnsafeQueue(t *testing.T) {
	unsafeQueue := NewFIFOQueue(5)
	fmt.Println("Test Operations Unsafe Queue")
	testQueueOperations(unsafeQueue, t)
}

func TestOpsCondQueue(t *testing.T) {
	condQueue := NewCondQueue(5)
	fmt.Println("Test Operations Conditional Queue")
	testQueueOperations(condQueue, t)
}

func BenchmarkUnsafeQueue(b *testing.B) {
	unsafeQueue := NewFIFOQueue(10000)
	for i := 0; i < b.N; i++ {
		benchQueueOperations(unsafeQueue)
	}
}

const (
	NumberOfGoroutines = 4
	NumberOfOperations = 10000
)

const (
	Len = iota
	Enqueue
	Dequeue
	Empty
)

func testConcurrentQueueAccess(queue Queue) {
	rand.Seed(time.Now().Unix())
	var wg sync.WaitGroup
	wg.Add(NumberOfGoroutines)

	for i := 0; i < NumberOfGoroutines; i++ {
		// Need to ensure that there are always enough values queued
		if i%2 == 0 {
			go func(i int) {
				defer wg.Done()
				for j := 0; j < NumberOfOperations; j++ {
					op := rand.Intn(4)
					switch op {
					case Len:
						queue.Len()
					case Enqueue:
						queue.Enqueue("Data" + strconv.Itoa(i) + strconv.Itoa(j))
					case Dequeue:
						queue.Dequeue()
					case Empty:
						queue.Empty()
					}

				}
			}(i)
		} else {
			go func(i int) {
				defer wg.Done()
				for j := 0; j < NumberOfOperations; j++ {
					queue.Enqueue(i)
				}
			}(i)
		}
	}
	wg.Wait()
}

func testQueueOperations(queue Queue, t *testing.T) {
	var length int
	var wantLength int
	var testName string
	var empty bool

	// Initial Queue Test
	testName = "Initial Queue Test"
	wantLength = 0
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	// Enqueue One Test
	testName = "Enqueue One Test"
	queue.Enqueue("Item1")
	wantLength = 1
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	item1 := queue.Dequeue()
	if item1 != "Item1" {
		t.Errorf("\n%s\nAction: Dequeue()\nWant: Item1\nGot: %v", testName, item1)
	}
	wantLength = 0
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	// Enqueue Three Test
	testName = "Enqueue Three Test"
	queue.Enqueue("Item2")
	queue.Enqueue(3)
	queue.Enqueue(4.0001)
	wantLength = 3
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	item2 := queue.Dequeue()
	if item2 != "Item2" {
		t.Errorf("\n%s\nAction: Dequeue()\nWant: Item2\nGot: %v", testName, item2)
	}
	wantLength = 2
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	item3 := queue.Dequeue()
	if item3 != 3 {
		t.Errorf("\n%s\nAction: Dequeue()\nWant: 3\nGot: %v", testName, item3)
	}
	wantLength = 1
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	item4 := queue.Dequeue()
	if item4 != 4.0001 {
		t.Errorf("\n%s\nAction: Dequeue()\nWant: 4.0001\nGot: %v", testName, item4)
	}
	wantLength = 0
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	wantLength = 0
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	// Queue Slice Allocation Test
	testName = "Queue Slice Allocation Test"
	size := 200
	for i := 0; i < size; i++ {
		queue.Enqueue(i)
	}
	wantLength = size
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	for j := 0; j < size; j++ {
		if x := queue.Dequeue(); x != j {
			t.Errorf("\n%s\nAction: Enqueue()\nWant: %d\nGot: %v", testName, j, x)
			break
		}
	}

	// Queue Flush Test
	testName = "Queue Flush Test"
	size = 200
	for i := 0; i < size; i++ {
		queue.Enqueue(i)
	}

	queue.Flush()
	wantLength = 0
	if length = queue.Len(); length != wantLength {
		t.Errorf("\n%s\nAction: Len()\nWant: %d\nGot: %d", testName, wantLength, length)
	}

	// Queue Empty Test
	testName = "Queue Empty Test"
	queue.Enqueue(1)
	queue.Enqueue("test")

	if empty = queue.Empty(); empty {
		t.Errorf("\n%s\nAction: Empty()\nWant: %t\nGot: %t", testName, false, empty)
	}

	queue.Dequeue()
	queue.Dequeue()
	if empty = queue.Empty(); !empty {
		t.Errorf("\n%s\nAction: Empty()\nWant: %t\nGot: %t", testName, true, empty)
	}
}

func benchQueueOperations(queue Queue) {
	const nrOfOps = 10000

	for i := 0; i < nrOfOps; i++ {
		queue.Enqueue(i)
	}

	for j := 0; j < nrOfOps; j++ {
		queue.Dequeue()
	}
}
