package all

import (
	"encoding/json"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Pool struct {
	tasks map[string]*Task
	mutex sync.Mutex
}

func NewPool() *Pool {
	return &Pool{tasks: make(map[string]*Task), mutex: sync.Mutex{}}
}

func (p *Pool) MarkTaskAsProcessing(id string) {
	p.mutex.Lock()
	_, ok := p.tasks[id]
	if ok {
		p.tasks[id].Status = WorkStatus
		p.tasks[id].StartTime = time.Now()
	}
	p.mutex.Unlock()
}

func (p *Pool) MarkTaskAsDone(id string, data string) {
	p.mutex.Lock()
	_, ok := p.tasks[id]
	if ok {
		p.tasks[id].Status = StopStatus
		p.tasks[id].StopTime = time.Now()
	}
	p.mutex.Unlock()
}

func (p *Pool) GetData(id string) interface{} {
	p.mutex.Lock()
	_, ok := p.tasks[id]
	var result interface{} = nil
	if ok {
		result = p.tasks[id].GetData()
	}

	p.mutex.Unlock()
	return result
}

func (p *Pool) GetStatus(id string) ([]byte, error, bool) {
	p.mutex.Lock()
	var data []byte
	var err error
	_, ok := p.tasks[id]
	if ok {
		p.tasks[id].SetStatus()
		data, err = json.Marshal(p.tasks[id])
	}

	p.mutex.Unlock()
	return data, err, ok
}

func (p *Pool) DeleteTask(id string) bool {
	p.mutex.Lock()
	_, ok := p.tasks[id]
	if ok {
		p.tasks[id].Stop <- struct{}{}
		delete(p.tasks, id)
	}
	p.mutex.Unlock()

	return ok
}

func (p *Pool) CreateTask() (JobDataToHandle, bool) {
	p.mutex.Lock()
	stop := make(chan struct{})
	result := JobDataToHandle{}
	id := uuid.New().String()
	_, ok := p.tasks[id]
	if !ok {
		data := "Some data from task "
		p.tasks[id] = NewTask(id, data, stop)
		result.Id = id
		result.InputData = data
		result.Stop = stop
	}

	p.mutex.Unlock()
	return result, !ok
}
