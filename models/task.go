package models

import (
	"time"
)

type Task[T any] struct {
	Id         string        `json:"-"`
	Data       T             `json:"-"`
	DateCreate string        `json:"dateCreate"`
	StartTime  time.Time     `json:"-"`
	StopTime   time.Time     `json:"-"`
	WorkTime   string        `json:"worktime"`
	Status     string        `json:"status"`
	Stop       chan struct{} `json:"-"`
}

func NewTask(id string, data interface{}, stop chan struct{}) *Task[any] {
	return &Task[any]{
		Id:         id,
		Status:     CreateStatus,
		DateCreate: time.Now().Format("2006-01-02 15:04:05"),
		Data:       data,
		Stop:       stop,
	}
}
