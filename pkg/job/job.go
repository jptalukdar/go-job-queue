package job

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Status string

const (
	StatusStart      Status = "START"
	StatusReady      Status = "READY"
	StatusRunning    Status = "RUNNING"
	StatusTerminated Status = "TERMINATED"
	StatusWait       Status = "WAITING"
	StatusInit       Status = "INITIALIZING"
)

type Response struct {
}

type JobMeta struct {
	Timestamp string
}

type Job struct {
	Id             string
	State          Status
	ExecutorMethod func(context.Context) Response
	JobMeta        JobMeta
	Ctx            context.Context
	Heatbeat       chan int
	StateComm      chan Status
	wg             sync.WaitGroup
}

func NewJob(f func(context.Context) Response) Job {
	j := Job{}
	j.State = StatusInit
	j.ExecutorMethod = f
	j.Ctx = context.Background()
	// j.JobMeta.Timestamp =
	return j
}
func (j *Job) Monitor() {

}

func (j *Job) Start() {
	if j.State == StatusRunning {
		log.Printf("Job already running")
		return
	}
	// Start the job
	log.Printf("Starting")
	j.wg.Add(1)

	// Trigger a goroutine while executing the method
	go func() {
		defer j.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		log.Printf("Job started")
		j.State = StatusRunning
		j.ExecutorMethod(j.Ctx)
		// j.StateComm <- StatusTerminated
		// log.Printf("Block")
		// j.Ctx.Done()
		log.Printf("Done")
	}()

	log.Printf("Started")
}

func (j *Job) StartWithWaitGroup(wg sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	j.Start()
	j.Wait()
}

func (j *Job) Stop() {
	// Stop the job
	j.Ctx.Done()
	j.State = StatusTerminated
}

func (j *Job) Check() string {
	// Check the status of job
	return string(j.State)

}

func (j *Job) Wait() {
	// Waits for job to finish
	// // Wait for it to complete
	j.wg.Wait()
}
