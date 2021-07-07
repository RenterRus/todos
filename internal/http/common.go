package http

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net/url"
	"strings"
	"time"
)

func (s *HTTPServer) agregator(form url.Values) map[string]string {
	httpArgs := map[string]string{}
	for k, v := range form {
		httpArgs[k] = strings.Join(v, "")
	}
	return httpArgs
}

func GetGRPCConnect(grpcAddr string) *grpc.ClientConn {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                60 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}))
	if err != nil {
		fmt.Println(err.Error())
	}
	return conn
}