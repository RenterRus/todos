package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"todoes/internal/sqlite"
	"todoes/pb/todoes"
)

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