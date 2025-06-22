package all

import (
	"time"
)

type Task struct {
	Id         string        `json:"-"`
	Data       interface{}   `json:"-"`
	DateCreate string        `json:"dateCreate"`
	StartTime  time.Time     `json:"-"`
	StopTime   time.Time     `json:"-"`
	WorkTime   string        `json:"worktime"`
	Status     string        `json:"status"`
	Stop       chan struct{} `json:"-"`
}

func NewTask(id string, data interface{}, stop chan struct{}) *Task {
	return &Task{
		Id:         id,
		Status:     CreateStatus,
		DateCreate: time.Now().Format("2006-01-02 15:04:05"),
		Data:       data,
		Stop:       stop,
	}
}

func (task *Task) SetStatus() {
	switch task.Status {
	case CreateStatus:
		task.WorkTime = (time.Millisecond * 0).String()
	case WorkStatus:
		task.WorkTime = time.Since(task.StartTime).String()
	case StopStatus:
		task.WorkTime = task.StopTime.Sub(task.StartTime).String()
	default:
		panic("Invalid status")
	}
}

func (task *Task) GetData() interface{} {
	switch task.Status {
	case CreateStatus:
		return "Данные отсутствуют"
	case WorkStatus:
		return "Данные отсутствуют"
	case StopStatus:
		return task.Data
	default:
		panic("Invalid status")
	}
}
