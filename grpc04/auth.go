package grpc04

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

var (
	userName = "admin"
	password = "admin"
)

type AuthService struct {
	Username string
	password string
}

// grpc 实现自定义身份验证
//type PerRPCCredentials interface {
//	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
//	RequireTransportSecurity() bool
//}

// GetRequestMetadata 定义授权信息的具体存放形式，最终会按这个格式存放到 metadata map 中。
func (a *AuthService) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	authData := map[string]string{
		"username": a.Username,
		"password": a.password,
	}
	return authData, nil
}

// RequireTransportSecurity 是否需要基于 TLS 加密连接进行安全传输
func (a *AuthService) RequireTransportSecurity() bool {
	return false
}

// NewAuthService 实例花AuthService
func NewAuthService() *AuthService {
	return &AuthService{
		Username: userName,
		password: password,
	}
}

// IsValidAuth 验证逻辑账号密码是否正确
func IsValidAuth(ctx context.Context) error {

	// 获取metadata 信息
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	// 获取metadata的账号密码信息
	if userName != md["username"][0] || password != md["password"][0] {
		return status.Errorf(codes.Unauthenticated, "Unauthorized")
	}
	log.Println("login success")
	return nil
}
