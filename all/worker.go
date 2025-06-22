package all

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Run(job JobDataToHandle) string {
	i := rand.Intn(3) + 3

	select {
	case <-job.Stop:
		return ""
	case <-time.After(time.Minute * time.Duration(i)):
		return job.InputData + job.Id
	}
}

type (
	JobStorage interface {
		MarkTaskAsProcessing(id string)
		MarkTaskAsDone(id string, data string)
	}
)

type (
	JobDataToHandle struct {
		Id        string
		InputData string
		Stop      chan struct{}
	}
)

type (
	Worker struct {
		jobs         <-chan JobDataToHandle
		stop         chan struct{}
		wg           *sync.WaitGroup
		countWorkers int
		jobStorage   JobStorage
		jobHandler   func(JobDataToHandle) string
	}
)

func (p *Worker) Close() {
	close(p.stop)
	p.wg.Wait()
}

func NewWorker(jobs <-chan JobDataToHandle, count int, js JobStorage) *Worker {
	w := &Worker{jobs: jobs, stop: make(chan struct{}), wg: &sync.WaitGroup{}, countWorkers: count, jobHandler: Run, jobStorage: js}
	for i := 0; i < w.countWorkers; i++ {
		w.wg.Add(1)
		go w.worker(i)
	}

	return w
}

func (p *Worker) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case j, ok := <-p.jobs:
			if !ok {
				return
			}
			fmt.Println("Worker", id, "processing")
			p.jobStorage.MarkTaskAsProcessing(j.Id)
			data := p.jobHandler(j)
			p.jobStorage.MarkTaskAsDone(j.Id, data)
			fmt.Println("Worker", id, "stop")
		case <-p.stop:
			return
		}
	}
}
