package lab5

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

type schedule struct {
	in          []Job
	want        []int
	policy      int
	quantum     time.Duration
	nrOfWorkers int
}

var scheduleTests = []schedule{
	{
		[]Job{},
		[]int{},
		FIFO,
		time.Second,
		1,
	},
	{
		[]Job{NewJob(0, nil, 20*time.Millisecond), NewJob(1, nil, 10*time.Millisecond)},
		[]int{0, 1},
		FIFO,
		time.Second,
		1,
	},
	{
		[]Job{NewJob(0, nil, 10*time.Millisecond), NewJob(1, nil, 20*time.Millisecond)},
		[]int{0, 1},
		FIFO,
		time.Second,
		1,
	},
	{
		[]Job{NewJob(0, nil, 30*time.Millisecond), NewJob(1, nil, 20*time.Millisecond), NewJob(2, nil, 40*time.Millisecond)},
		[]int{0, 1, 2},
		FIFO,
		time.Second,
		1,
	},
	{
		[]Job{NewJob(0, nil, 10*time.Millisecond), NewJob(1, nil, 20*time.Millisecond)},
		[]int{0, 1},
		SJF,
		time.Second,
		1,
	},
	{
		[]Job{NewJob(0, nil, 20*time.Millisecond), NewJob(1, nil, 10*time.Millisecond)},
		[]int{1, 0},
		SJF,
		time.Second,
		1,
	},
	{
		[]Job{NewJob(0, nil, 30*time.Millisecond), NewJob(1, nil, 10*time.Millisecond), NewJob(2, nil, 5*time.Millisecond), NewJob(3, nil, 20*time.Millisecond)},
		[]int{2, 1, 3, 0},
		SJF,
		time.Second,
		1,
	},
	{
		[]Job{NewJob(0, nil, 20*time.Millisecond)},
		[]int{0, 0, 0, 0},
		RR,
		5 * time.Millisecond,
		1,
	},
	{
		[]Job{NewJob(0, nil, 20*time.Millisecond), NewJob(1, nil, 30*time.Millisecond)},
		[]int{0, 1, 0, 1, 0, 1, 0, 1, 1, 1},
		RR,
		5 * time.Millisecond,
		1,
	},
	{
		[]Job{NewJob(0, nil, 40*time.Millisecond), NewJob(1, nil, 10*time.Millisecond), NewJob(2, nil, 5*time.Millisecond), NewJob(3, nil, 20*time.Millisecond)},
		[]int{0, 1, 2, 3, 0, 3, 0, 0},
		RR,
		10 * time.Millisecond,
		1,
	},
}

var multiWorkerSchedule = []schedule{
	{
		[]Job{NewJob(0, nil, 200*time.Millisecond), NewJob(1, nil, 45*time.Millisecond), NewJob(2, nil, 5*time.Millisecond), NewJob(3, nil, 30*time.Millisecond)},
		[]int{1, 2, 3, 0},
		FIFO,
		10 * time.Millisecond,
		2,
	},
	{
		[]Job{NewJob(0, nil, 200*time.Millisecond), NewJob(1, nil, 45*time.Millisecond), NewJob(2, nil, 5*time.Millisecond), NewJob(3, nil, 30*time.Millisecond)},
		[]int{2, 3, 1, 0},
		SJF,
		10 * time.Millisecond,
		2,
	},
	{
		[]Job{NewJob(0, nil, 45*time.Millisecond), NewJob(1, nil, 5*time.Millisecond), NewJob(2, nil, 30*time.Millisecond)},
		[]int{1, 0, 2, 0, 2, 0},
		RR,
		15 * time.Millisecond,
		2,
	},
	{
		[]Job{NewJob(0, nil, 200*time.Millisecond), NewJob(1, nil, 45*time.Millisecond), NewJob(2, nil, 5*time.Millisecond), NewJob(3, nil, 30*time.Millisecond)},
		[]int{2, 3, 1, 0},
		FIFO,
		10 * time.Millisecond,
		4,
	},
}

func TestScheduling(t *testing.T) {
	submitSchecule(t, scheduleTests)
}

func TestMultiWorker(t *testing.T) {
	submitSchecule(t, multiWorkerSchedule)
}

func BenchmarkSchedulingSimulate(b *testing.B) {
	const nrOfJobs = 500
	runtime.GOMAXPROCS(runtime.NumCPU())
	nrOfWorkers := runtime.NumCPU()
	jobs := make([]Job, nrOfJobs)
	for i := 0; i < nrOfJobs; i++ {
		jobs[i] = NewJob(i, simulateWork, time.Duration(rand.Intn(5))*time.Millisecond)
	}
	s := NewScheduler()
	go s.Schedule(jobs, FIFO, 300*time.Microsecond)
	s.CreateWorkerPool(nrOfWorkers)

	for len(s.Results) > 0 {
		<-s.Results
	}
}

func BenchmarkSchedulingSleep(b *testing.B) {
	const nrOfJobs = 500
	runtime.GOMAXPROCS(runtime.NumCPU())
	nrOfWorkers := runtime.NumCPU()
	jobs := make([]Job, nrOfJobs)
	for i := 0; i < nrOfJobs; i++ {
		jobs[i] = NewJob(i, nil, time.Duration(rand.Intn(5))*time.Millisecond)
	}
	s := NewScheduler()
	go s.Schedule(jobs, FIFO, 300*time.Microsecond)
	s.CreateWorkerPool(nrOfWorkers)

	for len(s.Results) > 0 {
		<-s.Results
	}
}

func simulateWork(d time.Duration) {
	start := time.Now()
	for time.Since(start) < d {

	}
}

func submitSchecule(t *testing.T, input []schedule) {
	for i, st := range input {
		s := NewScheduler()
		go s.Schedule(st.in, st.policy, st.quantum)
		s.CreateWorkerPool(st.nrOfWorkers)
		j := 0
		for res := range s.Results {
			if j >= len(st.want) {
				t.Errorf("schedule test %d, job %d: too many jobs, want %d jobs", i, j, len(st.want))
				break
			}
			if res.job.id != st.want[j] {
				t.Errorf("schedule test %d, job %d: got id %v for input %v, want %v", i, j, res.job.id, st.in, st.want[j])
			}
			j++
		}
	}
}
