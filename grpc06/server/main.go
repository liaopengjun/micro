package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"micro/grpc02/service"
	"net"
	"time"
)

type server struct {
	service.UnimplementedEchoServer
}

// UnaryEcho 一个普通的UnaryAPI
func (e *server) UnaryEcho(ctx context.Context, req *service.EchoRequest) (*service.EchoResponse, error) {
	message := req.Message
	if message == "Normal" {
		time.Sleep(time.Millisecond * 500)
		message = "hello" + message
	}
	if message == "Timeout" {
		time.Sleep(time.Second * 2)
		message = "hello" + message
	}
	return &service.EchoResponse{
		Message: message,
	}, nil
}

func main() {
	c := grpc.NewServer()
	service.RegisterEchoServer(c, new(server))
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err = c.Serve(listen); err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("failed to serve: %v", err)
	}
}
