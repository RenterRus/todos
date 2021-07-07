package grpc

import (
	"google.golang.org/grpc"
	"todoes/pb/todoes"
)

type GrpcServer struct {
	todoes.UnimplementedTodoesServer
	server *grpc.Server
}
