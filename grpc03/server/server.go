package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"micro/grpc03/service"
	"net"
	"time"
)

type server struct {
	service.UnimplementedEchoServer
}

// logger 简单打印日志
func logger(format string, a ...interface{}) {
	fmt.Printf("LOG:\t"+format+"\n", a...)
}

func (s *server) UnaryEcho(ctx context.Context, in *service.EchoRequest) (*service.EchoResponse, error) {
	fmt.Printf("unary echoing message %q\n", in.Message)
	return &service.EchoResponse{Message: in.Message}, nil
}

// unaryInterceptor 一元拦截器：记录请求日志
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	m, err := handler(ctx, req)
	end := time.Now()
	// 记录请求参数 耗时 错误信息等数据
	logger("RPC: %s,req:%v start time: %s, end time: %s, err: %v", info.FullMethod, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return m, err
}

func main() {

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("failed to listen:", err)
		return
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	service.RegisterEchoServer(s, new(server))
	s.Serve(lis)
}
