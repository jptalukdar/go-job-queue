package main

import (
	"context"
	"fmt"
	"time"

	pkg "github.com/jptalukdar/go-job-queue/pkg/job"
	"github.com/jptalukdar/go-job-queue/pkg/queue"
)

func RunOsSleepCommand(ctx context.Context) pkg.Response {
	fmt.Println("I am a command running")
	// cmd := exec.Command("go")

	// err := cmd.Run()

	// if err != nil {
	// 	log.Fatal(err)
	// 	// panic(err)
	// }
	return pkg.Response{}
}

func main() {
	job := pkg.NewJob(RunOsSleepCommand)

	q := queue.NewQueue()
	q.Run()

	q.Add(job)
	// job.Start()
	// job.Wait()

	time.Sleep(5 * time.Second)

	j := pkg.NewJob(RunOsSleepCommand)
	q.Add(j)

	q.Wait()
	// j.Start()
	// j.Wait()
}
