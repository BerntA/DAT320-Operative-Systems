![UiS](https://www.uis.no/getfile.php/13391907/Biblioteket/Logo%20og%20veiledninger/UiS_liggende_logo_liten.png)

# Lab 4: Data Race Detection & Profiling

| Lab 4:		| Data Race Detection & Profiling		|
| -------------------- 	| ------------------------------------- |
| Subject: 		| DAT320 Operating Systems 		|
| Deadline:		| Oct 12th 2018 16:00			|
| Expected effort:	| 10-15 hours 				|
| Grading: 		| Pass/fail 				|
| Submission: 		| Individually				|


## Introduction

This lab exercise is divided into two parts and deals with two separate programming tools. The
first part will focus on high-level synchronization techniques and will give an introduction to
Go’s built-in data race detector. You will use two different techniques to ensure synchronization
to a shared data structure. The second part of the lab deals with CPU and memory profiling.
We will analyze different implementations of a simple data structure.


## High-level Synchronization & Data Race Detection

In the first part of this lab we will focus on high-level synchronization techniques using the Go
programming language. We will return to these topics also in lab 5.

The Go language provides a built-in race detector that we will use to identify data races and
verify implementations. A data race occurs when two threads (or goroutines) access a variable
concurrently and at least one of these accesses is a write.

We will work on a stack data structure that will be accessed concurrently from several
goroutines. The stack stores values of type interface{}, meaning any value type. The stack
interface is shown in Listing 1 and is found in the file common.go. The interface contains three
methods. Len() returns the current number of items on the stack, Pop() interface{} pops
an item of the stack (nil if empty), while Push(value interface{}) pushes an item onto the
stack.

Listing 1: Stack interface
```
type Stack interface {
	Len() int
	Push(value interface{})
	Pop() interface{}
}
```

For this lab we will use the tests defined in stacks test.go to verify the different stack
implementations. The tests can be run using the go test command. We will run one test at a
time. Running only a specific test can be achieved by supplying the -run flag together with a
regular expression indicating the test names. For example, to run only the TestUnsafeStack
function, use go test -run TestUnsafeStack. There are two type of tests defined for each
stack implementation we will be working on. One test verifies a stack’s operations, while the
other is meant to test concurrent access using a race detector. Study the test file for details.

As stated in the introduction, Go includes a built-in data race detector. Read [Data Race Detector]
(http://golang.org/doc/articles/race_detector.html) for an introduction and usage examples.

**Race Questions**
```
1.) What do you call a piece of code that must only be accessed by one process at a time?
	a. monitor
	b. race condition
	c. critical section
	d. semaphore

2.) Which command will run the race detector along with the test TestUnsafeStack?
	a. go test -race TestUnsafeStack
	b. go test -race -run TestUnsafeStack
	c. go test -run TestUnsafeStack
	d. go test -run -race TestUnsafeStack

3.) Which method(s) in UnsafeStack can cause a race condition?
	a. Len()
	b. Push()
	c. Pop()
	d. All of the above
```

**Tasks**
```
1. The file stack_sync.go is a copy of stack.go, but the type is renamed to SafeStack.
Modify this file so that access to the SafeStack type is synchronized (can be accessed
safely from concurrently running goroutines). You can use the Mutex type from the sync
package to achieve this.

2. Verify your implementation by running the TestSafeStack test with the data race detector
enabled. The test should not produce any data race warnings.

go test -race -run TestSafeStack

3. Go has a built-in high-level API for concurrent programming based on Communicating
Sequential Processes (CSP). This API promotes synchronization through sending and receiving 
data via thread-safe channels (as opposed to traditional locking). 
The file stack csp.go contains a CspStack type that implements the stack interface in
Listing 1 (but the actual method implementations are empty). The type also has a constructor 
function needed for this task. Modify this file so that access to the CspStack type
is synchronized. The synchronization should be achieved by using Go’s CSP features (channels 
and goroutines). 
```
There is in this case an amount of overhead when using channels to achieve synchronization 
compared to locking. The main point for this task is to give an introduction on how to use 
channels (CSP) for synchronization. 
This will require some self-study if you are not familiar with Go’s
CSP-based concurrent programming capabilities. A place to start can be the introduction
found [here](http://golang.org/doc/effective_go.html#concurrency). For inspiration on
how to solve this specific task, look at Chapter 7.2.3 in M. Summerfield, Programming in
Go: Creating Applications for the 21st Century (Developer’s Library), 1st ed. Addison-
Wesley Professional, 5 2012. The sub chapter is available as a [pdf](https://stavanger.instructure.com/files/67324/download?download_frd=1) on canvas. Note
that you should also ensure that the stack operations are implemented correctly. You can
verify them by running the TestOpsCspStack test.
```
4. Verify your implementation by running the TestCspStack test with the data race detector
enabled. The test should not produce any data race warnings.

go test -race -run TestCspStack

5. Optional exercise: Implement a ForcePop() method for the SafeStack type or the
CspStack. The method should pop an item off the stack, but if the stack is empty, it
should block until an item is pushed onto it. For the SafeStack type an option may be to
use the Cond type in the sync package. To use channels would be a natural choice if you
choose to adjust the CSP-version.
```

**Recommended reading**

[The Go Memory Model](http://golang.org/ref/mem)

[Introducing the Go Race Detector](http://blog.golang.org/race-detector)


## CPU and Memory Profiling

In this part of the lab we will use a technique called profiling to dynamically analyze a program.
Profiling can among other things be used to measure an application’s CPU utilization and
memory usage. Being able to profile applications is very helpful for doing optimizations and
is an important part of Systems Programming. This lab will give a very short introduction to
how profiling data can be analyzed. You may in future lab exercises be required to use profiling
to improve and optimize programs.

Profiling for Go can be enabled through the runtime/pprof package or by using the testing
package’s profiling support. Profiles can be analyzed and visualized using the go tool pprof
program.

We will continue to use the stack implementations used in the first part of the lab. The file
stacks test.go contains one benchmark for the three different implementations. Each of
them uses the same core stack benchmark defined in the benchStackOperations(stack Stack)
function. The stack implementations are not accessed concurrently so that the benchmarks can
be kept reasonably deterministic.

Read [Profiling Go Programs](http://blog.golang.org/profiling-go-programs). This
blog post present a good introduction to Go’s profiling abilities. You should also look at 
[testing](http://golang.org/pkg/testing/) and [testing flags](http://golang.org/cmd/go/#Description_of_testing_flags) 
for information on how to run the benchmarks and details about how Go’s
testing tool easily enables profiling when benchmarking.


**Profiling Questions**
```
1.) Which testing flag changes the amount of memory allocations measured?
	a. -benchmem
	b. -memprofile
	c. -memprofilerate
	d. -benchtime

2.) Which pprof command will list the top 20 profiles?
	a. topN
	b. top20
	c. top 20
	d. 20 top

3.) Which pprof command will produce a graph of the profile and show it on a web browser?
	a. web
	b. list
	c. svg
	d. top

4.) Which stack implementation uses the most CPU and memory?
	a. SafeStack
	b. CspStack
	c. SliceStack
```

**Tasks**
```
1. The file stack slice.go contains a stack implementation, SliceStack, backed by a slice
(dynamic array). You will need adjust this implementation to be synchronized in the exact
same way you did for the SafeStack type. This has to be done to make the benchmark
between the three implementations fair and comparable.

2. Run the three stack benchmarks described earlier. Enable memory allocation statistics
by supplying the -benchmem flag. In addition also use -memprofilerate=1. 
Attach the benchmark output in your report. You can avoid running any tests together 
with the benchmarks by providing a non-matching regular expression to the run flag 
(for example -run none).

3. Run the BenchmarkCspStack separately. This time you should also write a CPU profile to
file when running it. Load the data in the pprof program. Attach the top ten and top
ten cumulative listing of function samples to your report. Look at the top ten cumulative
listing.

4. Run the BenchmarkSafeStack separately and write a memory profile to file. Using pprof
identify the only function allocating memory. List this function and identify the line number
where the allocations occur. Also attach the output of this listing to your report.

5. Optional exercise: Explore the visualization possibilities offered by go tool pprof when
analyzing profiling data.
```

## Lab Approval

To have your lab assignment approved, you must come to the lab during lab hours
and present your solution. This lets you present the thought process behind
your solution, and gives us more information for grading purposes. When you are
ready to show your solution, reach out to a member of the teaching staff.  It
is expected that you can explain your code and show how it works. You may show
your solution on a lab workstation or your own computer. The results from
Autograder will also be taken into consideration when approving a lab. At least
60% of the Autograder tests should pass for the lab to be approved. A lab needs
to be approved before Autograder will provide feedback on the next lab
assignment.

Also see the [Grading and Collaboration Policy](https://github.com/uis-dat320-fall18/course-info/blob/master/policy.md)
