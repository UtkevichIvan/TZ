package all

import (
	"fmt"
	"net/http"
)

type MyServer struct {
	jobs         chan JobDataToHandle
	wp           *Worker
	mux          *http.ServeMux
	tasks        *Pool
	port         string
	countWorkers int
	sender       *Sender
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
	data, err, ok := s.tasks.GetStatus(id)

	if !ok {
		http.Error(w, "There is no task with such ID", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *MyServer) GetResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data := s.tasks.GetData(id)
	if data == nil {
		http.Error(w, "There is no task with such ID", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, data)
}

func NewMyServer() *MyServer {
	j := make(chan JobDataToHandle)
	cW := 4
	tasks := NewPool()
	return &MyServer{
		jobs:         j,
		wp:           NewWorker(j, cW, tasks),
		mux:          http.NewServeMux(),
		tasks:        tasks,
		port:         "8080",
		countWorkers: cW,
		sender:       NewSender(j),
	}
}

func (s *MyServer) Start() {
	go s.sender.Send()
	s.mux.HandleFunc("/create", s.CreateTask)
	s.mux.HandleFunc("/delete/{id}", s.DeleteTask)
	s.mux.HandleFunc("/status/{id}", s.GetStatus)
	s.mux.HandleFunc("/data/{id}", s.GetResult)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", s.mux)
}
