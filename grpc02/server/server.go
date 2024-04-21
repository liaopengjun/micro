package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"micro/grpc02/service"
	"net"
)

// UnaryAPI：普通一元方法
type Echo struct {
	service.UnimplementedEchoServer
}

// UnaryEcho 一个普通的UnaryAPI

func (e *Echo) UnaryEcho(ctx context.Context, req *service.EchoRequest) (*service.EchoResponse, error) {
	log.Printf("UnaryEcho Recved: %v", req.GetMessage())
	resp := &service.EchoResponse{
		Message: req.GetMessage(),
	}
	return resp, nil
}

// ServerStreamingEcho 客户端发送一个请求 服务端以流的形式循环发送多个响应
/*
1. 获取客户端请求参数
2. 处理完成后返回过个响应
3. 最后返回nil表示已经完成响应
*/
func (e *Echo) ServerStreamingEcho(req *service.EchoRequest, stream service.Echo_ServerStreamingEchoServer) error {
	log.Printf("ServerStreaming Recved: %v", req.GetMessage())
	// 接收到一个hello world 循环返回10个hello world
	for i := 0; i < 10; i++ {
		err := stream.Send(&service.EchoResponse{
			Message: req.GetMessage(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// 构造一个grpc对象
	conn := grpc.NewServer()
	// 注册函数
	service.RegisterEchoServer(conn, new(Echo))
	// 监听端口
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Printf("listen error: %v", err)
		return
	}
	// 启动grpc服务
	conn.Serve(lis)
}
