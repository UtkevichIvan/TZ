package controller

import (
	"TZ/models"
	"TZ/worker"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type MyServer struct {
	jobs         chan all.JobDataToHandle
	wp           *all.Worker
	mux          *http.ServeMux
	tasks        *models.Pool
	port         string
	countWorkers int
	sender       *all.Sender
}

func (s *MyServer) CreateTask(w http.ResponseWriter, r *http.Request) {
	data, ok := s.tasks.CreateTask()
	if !ok {
		http.Error(w, "A task with this ID already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data.Id))
	s.sender.Add(data)
}

func (s *MyServer) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ok := s.tasks.DeleteTask(id)
	if !ok {
		http.Error(w, "There is no task with such ID", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *MyServer) GetStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, ok := s.tasks.GetStatus(id)

	if !ok {
		http.Error(w, "There is no task with such ID", http.StatusBadRequest)
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func (s *MyServer) GetData(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data := s.tasks.GetData(id)
	if data == nil {
		http.Error(w, "There is no task with such ID", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, data)
}

func NewMyServer() *MyServer {
	j := make(chan all.JobDataToHandle)
	tasks := models.NewPool()
	return &MyServer{
		jobs:   j,
		wp:     all.NewWorker(j, 4, tasks),
		mux:    http.NewServeMux(),
		tasks:  tasks,
		sender: all.NewSender(j),
	}
}

func (s *MyServer) Start() {
	go s.sender.Start()
	s.mux.HandleFunc("/create", s.CreateTask)
	s.mux.HandleFunc("/delete/{id}", s.DeleteTask)
	s.mux.HandleFunc("/status/{id}", s.GetStatus)
	s.mux.HandleFunc("/data/{id}", s.GetData)
	log.Println("Server is listening...")
	err := http.ListenAndServe(":8080", s.mux)
	if err != nil {
		log.Fatal(err)
		return
	}
}
