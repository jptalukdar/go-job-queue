package queue

import (
	"log"
	"sync"

	"github.com/jptalukdar/go-job-queue/pkg/job"
)

type channel chan job.Job

var lock = &sync.Mutex{}
var queueAddLock = &sync.Mutex{}
var queueInstance *queue
var channelInstance channel

type queue struct {
	queueChannel chan job.Job
	exitChannel  chan bool
	wg           sync.WaitGroup
	Running      bool
}

func (q *queue) Run() { // Should be called once
	if q.Running == true {
		return
	} else {
		q.Running = true
	}

	go func() {

		var wg sync.WaitGroup
		var Exit bool = false

		q.wg.Add(1) //Keep tracking of Server thread
		defer q.wg.Done()
		for {
			select {
			case j := <-q.queueChannel:
				j.StartWithWaitGroup(wg)
			case <-q.exitChannel:
				Exit = true
			default:

			}

			if Exit == true {
				break
			}
		}
		wg.Wait()
	}()
}

func (q *queue) Wait() {
	q.wg.Wait()
}
func (q *queue) Add(j job.Job) {
	queueAddLock.Lock()
	defer queueAddLock.Unlock()
	q.queueChannel <- j // Should be non blocking since buffered chan
}

func NewQueue() *queue {
	q := queue{}
	q.queueChannel = make(chan job.Job, 10) // Buffered channel. TODO: Set Buffer length
	return &q
}

func newChannel() channel {
	ch := make(chan job.Job)
	return ch
}

func getQueue() *queue {
	if queueInstance == nil {
		lock.Lock() //Locking is required only when creating a new queue
		defer lock.Unlock()
		// Checking is required twice if another instance has not already created one instance after acquring the lock
		if queueInstance == nil {
			queueInstance = NewQueue()
		}
	}
	return queueInstance
}

func getQueueChannel() channel {
	if channelInstance == nil {
		lock.Lock() //Locking is required only when creating a new queue
		defer lock.Unlock()
		// Checking is required twice if another instance has not already created one instance after acquring the lock
		if channelInstance == nil {
			channelInstance = newChannel()
		}
		log.Printf("Primary channel created")
	}
	return channelInstance
}

func Add(j job.Job) {
	q := getQueue()
	q.Add(j)
}

func init() {
	getQueueChannel()
	// getQueue() // Initializing when package is imported
}
