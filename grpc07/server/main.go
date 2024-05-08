package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"micro/grpc02/service"
	"net"
	"sync"
)

type failingServer struct {
	service.UnimplementedEchoServer
	mu         sync.Mutex
	reqCounter uint // 请求次数
	reqModulo  uint // 取模
}

// maybeFailRequest 手动模拟请求失败 一共请求n次，前n-1次都返回失败，最后一次返回成功。
func (s *failingServer) maybeFailRequest() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reqCounter++
	if (s.reqModulo > 0) && (s.reqCounter%s.reqModulo == 0) {
		return nil
	}
	return status.Errorf(codes.Unavailable, "maybeFailRequest: failing it")
}

func (s *failingServer) UnaryEcho(ctx context.Context, req *service.EchoRequest) (*service.EchoResponse, error) {
	if err := s.maybeFailRequest(); err != nil {
		log.Println("request failed count:", s.reqCounter)
		return nil, err
	}
	log.Println("request succeeded count:", s.reqCounter)
	return &service.EchoResponse{Message: req.Message}, nil
}

func main() {

	s := grpc.NewServer()

	service.RegisterEchoServer(s, &failingServer{
		reqModulo: 4,
	})

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err = s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("failed to serve: %v", err)
	}
}
