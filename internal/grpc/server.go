package grpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"todoes/internal/sqlite"
	"todoes/pb/todoes"
)

type GrpcServer struct {
	todoes.UnimplementedTodoesServer
	server *grpc.Server
}

func (s *GrpcServer) CreateTodo(ctx context.Context, request *todoes.CreateTodoRequest) (*todoes.CreateTodoResponse, error){
	deadline := &timestamp.Timestamp{
		Seconds: 0,
		Nanos: 0,
	}

	err := sqlite.DBClient.Write(request.Message, "open", *deadline)
	if err != nil {
		return nil, err
	}
	return &todoes.CreateTodoResponse{Code: 200, Msg: "Ok"}, nil
}
func (s *GrpcServer) GetTodos(ctx context.Context, request *empty.Empty) (*todoes.GetTodoResponse, error){
	res, err := sqlite.DBClient.ReadAll()
	if err != nil {
		return nil, err
	}
	return &todoes.GetTodoResponse{Code: 200, Msg: res}, nil
}

func (s *GrpcServer) DeleteTodo(ctx context.Context, request *todoes.DeleteTodoRequest) (*todoes.DeleteTodoResponse, error){
	err := sqlite.DBClient.CloseTodo(int(request.IdTodo))
	if err != nil{
		return &todoes.DeleteTodoResponse{Code: 500, Msg: err.Error()}, err
	}
	return &todoes.DeleteTodoResponse{Code: 200, Msg: "Ok"}, nil
}

func (s *GrpcServer) UpdateTodo(ctx context.Context, request *todoes.UpdateTodoRequest) (*todoes.UpdateTodoResponse, error){
	err := sqlite.DBClient.UpdateTodo(int(request.IdTodo), request.Message)
	if err != nil{
		return &todoes.UpdateTodoResponse{Code: 500, Msg: err.Error()}, err
	}
	return &todoes.UpdateTodoResponse{Code: 200, Msg: "Ok"}, nil
}


func (s *GrpcServer) SeTodoExpirationTimeoutTodo(ctx context.Context,
	request *todoes.SeTodoExpirationTimeoutTodoRequest) (*todoes.SeTodoExpirationTimeoutTodoResponse, error){
	fmt.Println("GRPC")
	err := sqlite.DBClient.SetExpTimeout(int(request.IdTodo), request.Deadline)
	if err != nil{
		return &todoes.SeTodoExpirationTimeoutTodoResponse{Code: 500, Msg: err.Error()}, err
	}
	return &todoes.SeTodoExpirationTimeoutTodoResponse{Code: 200, Msg: "Ok"}, nil
}

