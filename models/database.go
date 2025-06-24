package models

import (
	"TZ/worker"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

type Pool struct {
	tasks map[string]*Task[any]
	mutex sync.Mutex
}

func NewPool() *Pool {
	return &Pool{tasks: make(map[string]*Task[any]), mutex: sync.Mutex{}}
}

func (p *Pool) MarkTaskAsProcessing(id string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.tasks[id]
	if ok {
		p.tasks[id].Status = WorkStatus
		p.tasks[id].StartTime = time.Now()
	}
}

func (p *Pool) MarkTaskAsDone(id string, data string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.tasks[id]
	if ok {
		p.tasks[id].Data = data
		p.tasks[id].Status = StopStatus
		p.tasks[id].StopTime = time.Now()
	}
}

func (p *Pool) GetData(id string) interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.tasks[id]
	if ok {
		switch p.tasks[id].Status {
		case CreateStatus:
			return "Данные отсутствуют"
		case WorkStatus:
			return "Данные отсутствуют"
		case StopStatus:
			return p.tasks[id].Data
		}
	}

	return nil
}

func (p *Pool) GetStatus(id string) (*Task[any], bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.tasks[id]
	if ok {
		switch p.tasks[id].Status {
		case CreateStatus:
			p.tasks[id].WorkTime = (time.Millisecond * 0).String()
		case WorkStatus:
			p.tasks[id].WorkTime = time.Since(p.tasks[id].StartTime).String()
		case StopStatus:
			p.tasks[id].WorkTime = p.tasks[id].StopTime.Sub(p.tasks[id].StartTime).String()
		}
	}

	return p.tasks[id], ok
}

func (p *Pool) DeleteTask(id string) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.tasks[id]
	if ok {
		if p.tasks[id].Status != StopStatus {
			p.tasks[id].Stop <- struct{}{}
		}
		delete(p.tasks, id)
		log.Println("Task with ID: " + id + " has been deleted")
	}

	return ok
}

func (p *Pool) CreateTask() (all.JobDataToHandle, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	stop := make(chan struct{})
	result := all.JobDataToHandle{}
	id := uuid.New().String()
	_, ok := p.tasks[id]
	if !ok {
		data := "Some data from task "
		p.tasks[id] = NewTask(id, data, stop)
		result.Id = id
		result.InputData = data
		result.Stop = stop
		log.Println("Create task with id : " + id)
	}

	return result, !ok
}
