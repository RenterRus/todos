package grpc

import (
	"context"
	"fmt"
	"todoes/internal/sqlite"
	"todoes/pb/todoes"
)

func (s *GrpcServer) SeTodoExpirationTimeoutTodo(ctx context.Context,
	request *todoes.SeTodoExpirationTimeoutTodoRequest) (*todoes.SeTodoExpirationTimeoutTodoResponse, error){
	fmt.Println("GRPC")
	err := sqlite.DBClient.SetExpTimeout(int(request.IdTodo), request.Deadline)
	if err != nil{
		return &todoes.SeTodoExpirationTimeoutTodoResponse{Code: 500, Msg: err.Error()}, err
	}
	return &todoes.SeTodoExpirationTimeoutTodoResponse{Code: 200, Msg: "Ok"}, nil
}

