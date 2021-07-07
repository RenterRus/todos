package grpc

import (
	"context"
	"todoes/internal/sqlite"
	"todoes/pb/todoes"
)

func (s *GrpcServer) UpdateTodo(ctx context.Context, request *todoes.UpdateTodoRequest) (*todoes.UpdateTodoResponse, error){
	err := sqlite.DBClient.UpdateTodo(int(request.IdTodo), request.Message)
	if err != nil{
		return &todoes.UpdateTodoResponse{Code: 500, Msg: err.Error()}, err
	}
	return &todoes.UpdateTodoResponse{Code: 200, Msg: "Ok"}, nil
}