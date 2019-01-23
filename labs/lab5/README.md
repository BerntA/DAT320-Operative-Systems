![UiS](https://www.uis.no/getfile.php/13391907/Biblioteket/Logo%20og%20veiledninger/UiS_liggende_logo_liten.png)

# Lab 5: Threads and Protection

| Lab 5:		| Threads and Protection		|
| -------------------- 	| ------------------------------------- |
| Subject: 		| DAT320 Operating Systems 		|
| Deadline:		| Oct 26th 2018 16:00			|
| Expected effort:	| 10-15 hours 				|
| Grading: 		| Pass/fail 				|
| Submission: 		| Individually				|


### Table of Contents

1. [Introduction](https://github.com/uis-dat320-fall18/assignments/blob/master/lab5/README.md#introduction)
2. [Condition Variables](https://github.com/uis-dat320-fall18/assignments/blob/master/lab5/README.md#condition-variables)
3. [Job Scheduler](https://github.com/uis-dat320-fall18/assignments/blob/master/lab5/README.md#job-scheduler)
4. [Trace Tool](https://github.com/uis-dat320-fall18/assignments/blob/master/lab5/README.md#trace-tool)
5. [Lab Approval](https://github.com/uis-dat320-fall18/assignments/blob/master/lab5/README.md#lab-approval)


## Introduction

In this lab, you will work with condition variables in order to make a queue
implementation thread safe and prevent dequeueing when the queue is empty.
The other task revolves around building a job scheduler able to schedule jobs
according to different scheduling policies. Lastly, you will use the Trace
tool in Go to inspect the execution of the scheduler.

## Condition Variables

In this part of the lab, we will do a quick exercise using condition variables
in Go. As presented in the lectures, condition variables allow us to create
sections of code that are thread safe and only execute when some condition is met.
You can read more about condition variables in Go [here](https://golang.org/pkg/sync/#Cond).

**Task**

Provided in the `queues_fifo.go` file is a working implementation of a FIFO queue.
Your task is to create a wrapper around this queue that makes it thread safe.
The wrapper should call the methods implemented in the `FIFOQueue` type,
with locking to secure thread safety.

However, there is one detail that would make it difficult to implement using
regular mutexes; Dequeue should not return a value while the queue is empty.
In other words, the dequeue operation should wait for a new value to be enqueued.

The implementation should be written in `queues_cond.go` and should pass the tests
in `queues_test.go`.

## Job Scheduler

From the lectures you have learned about job scheduling. In this exercise, you will
be tasked with implementing a job scheduler in Go. This scheduler should be capable
of scheduling a list of jobs according to the different scheduling policies given in
the table below.

| Policy                    | Description                                                                                                                                                                                     |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| FIFO (First In First Out) | Schedules jobs in the order of arrival.                                                                                                                                                         |
| SJF (Shortest Job First)  | Schedules jobs based on the estimated execution time; runs the shortest jobs first.                                                                                                             |
| RR (Round Robin)          | Schedules jobs in the FIFO order, but only schedules the jobs for some give time quantum.<br> Therefore, not all jobs will be completed the first time, meaning the remaining time must be tracked. |


**Task**

The skeleton code for this task is located in `scheduler.go`, while the tests are in
`scheduler_test.go`. All the necessary types are already given. 

A `Job` object describes a job in the scheduler, and as such it keeps track of when it was
started, the estimated execution time, how much time is currently scheduled, and the remaining time.
It also references a function `task` that is used to simulate a job running.

When a job has been completed, a `Result` object should be sent, indicating that a job has been
executed. A result should be sent even if a job has not been completed, which using the
Round Robin policy. The `Result` object should contain a non-zero latency value when the job is
finished. Latency is the duration between the starting time of the job and when the job is completed.

Lastly, the `Scheduler` type is responsible for giving jobs to the workers through the
`Jobs` channel and outputting the results on the `Results` channel. Keep in mind that
the scheduler should be able to handle **at least 500 jobs.** It is also important to
close the two channels when all the data has been transmitted, such that the tests or
other parts of the program may stop waiting for more values over these channels.

The `Schedule` method should schedule the given jobs according to the provided policy.
You can safely ignore the `quantum` variable for all policies except for Round Robin.

`CreateWorkerPool` starts the given amount of workers, where a worker is implemented
in the `worker` method. Hint: You may find [WaitGroup](https://golang.org/pkg/sync/#WaitGroup) useful.

Finally, the `worker` method should be responsible for actually executing jobs.
A job should be executed by simply calling the `task` function with the correct duration.

## Trace Tool

In the previous lab you used the pprof tool to profile Go code. It is a good tool for
profiling CPU and memory usage, and finding problematic sections in the code.
Go has another profiling tool named Trace, focusing more on profiling concurrency
and latency in applications.

Below are some useful links about the tracing tool with examples:

* [Command Trace](https://golang.org/cmd/trace/)
* [go tool trace](https://making.pusher.com/go-tool-trace/)
* [An introduction to go tool trace](https://about.sourcegraph.com/go/an-introduction-to-go-tool-trace-rhys-hiltner/)
* [Using the Go tracer to speed up fractal making](https://campoy.cat/blog/using-the-go-tracer-to-speed-up-fractal-making/)

**Task**

Perform a trace while running the benchmarks seperately found in `scheduler_test.go` and inspect the
resulting trace files. Play around in the trace and inspect the features the tool has
to offer.

The two benchmarks are almost the same, the one difference being the task function for all the jobs.
`BenchmarkSchedulingSleep` uses the `time.Sleep` function, which leaves the threads idle for a while,
allowing you to see when the individual jobs are started. On the other hand, `BenchmarkSchedulingSimulate`
uses a function to keep the thread busy.

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

Also see the [Grading and Collaboration
Policy](https://github.com/uis-dat320-fall18/course-info/blob/master/policy.md)
document for additional information.
