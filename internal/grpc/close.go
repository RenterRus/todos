package grpc

import (
	"context"
	"todoes/internal/sqlite"
	"todoes/pb/todoes"
)

func (s *GrpcServer) DeleteTodo(ctx context.Context, request *todoes.DeleteTodoRequest) (*todoes.DeleteTodoResponse, error){
	err := sqlite.DBClient.CloseTodo(int(request.IdTodo))
	if err != nil{
		return &todoes.DeleteTodoResponse{Code: 500, Msg: err.Error()}, err
	}
	return &todoes.DeleteTodoResponse{Code: 200, Msg: "Ok"}, nil
}