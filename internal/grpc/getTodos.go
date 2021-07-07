package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"todoes/internal/sqlite"
	"todoes/pb/todoes"
)

func (s *GrpcServer) GetTodos(ctx context.Context, request *empty.Empty) (*todoes.GetTodoResponse, error){
	res, err := sqlite.DBClient.ReadAll()
	if err != nil {
		return nil, err
	}
	return &todoes.GetTodoResponse{Code: 200, Msg: res}, nil
}
