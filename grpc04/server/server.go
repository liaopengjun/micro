package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"micro/grpc02/service"
	"micro/grpc04"
	"net"
	"strings"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

type EcServer struct {
	service.UnimplementedEchoServer
}

func (e *EcServer) UnaryEcho(ctx context.Context, req *service.EchoRequest) (*service.EchoResponse, error) {
	log.Printf("UnaryEcho: %v", req.GetMessage())
	return &service.EchoResponse{Message: req.Message}, nil
}

// myEcUnaryInterceptor 中间件校验拦截器
func myEcUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	err = grpc04.IsValidAuth(ctx)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

// valid 校验认证信息有效性。
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "some-secret-token"
}

// ensureValidToken 用于校验 token 有效性的一元拦截器。
func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 如果 token不存在或者无效，直接返回错误，否则就调用真正的RPC方法。
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}

func main() {

	// 注册拦截器
	//s := grpc.NewServer(grpc.UnaryInterceptor(ensureValidToken)) // token 一元拦截器校验
	s := grpc.NewServer(grpc.UnaryInterceptor(myEcUnaryInterceptor)) // 账号密码一元拦截器校验
	// 注册服务
	service.RegisterEchoServer(s, new(EcServer))
	// 监听端口
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 启动 gRPC 服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
