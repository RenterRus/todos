package http

import (
	"context"
	"net/http"
	"todoes/pb/todoes"
)

type HTTPServer struct {
	s                  http.Server
	grpc  todoes.TodoesClient
}

type response struct {
	Err int32
	Messgae string
}

type Todo struct{
	Id int32 `json:"Id"`
	Message string `json:"Message"`
	Deadline int64 `json:"Deadline"`
	Status string `json:"Status"`
}

type ResponseTodo struct{
	Message []Todo
}

func NewServer(addr, grpcAddr string) *HTTPServer {
	s := new(HTTPServer)
	s.s = http.Server{
		Addr: addr,
	}
	s.grpc = todoes.NewTodoesClient(GetGRPCConnect(grpcAddr))
	return s
}

func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/GetTodos", s.GetTodos)
	mux.HandleFunc("/Create", s.CreateTodos)
	mux.HandleFunc("/Update", s.UpdateTodo)
	mux.HandleFunc("/Close", s.CloseTodos)
	mux.HandleFunc("/SetExpTimeout", s.SetExpTimeout)
	s.s.Handler = mux
	return s.s.ListenAndServe()
}

func (s *HTTPServer) GraceShutdown() {
	s.s.Shutdown(context.Background())
}
