package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"micro/grpc03/service"
	"micro/grpc04"
)

func unAry(client service.EchoClient) {
	resp, err := client.UnaryEcho(context.Background(), &service.EchoRequest{Message: "hello world"})
	if err != nil {
		fmt.Printf("client unary echo err:%v\n", err)
	}
	fmt.Printf("client unary echo resp:%v\n", resp)
}

// fetchToken 获取授权信息
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}

func main() {
	myAuth := grpc04.NewAuthService()

	// token 校验机制
	//perRPC := oauth.NewOauthAccess(fetchToken())
	//conn, err := grpc.Dial("127.0.0.1:1234", grpc.WithPerRPCCredentials(perRPC), grpc.WithInsecure())

	conn, err := grpc.Dial("127.0.0.1:1234", grpc.WithPerRPCCredentials(myAuth), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("client connect err:%v\n", err)
		return
	}
	defer conn.Close()

	client := service.NewEchoClient(conn)

	unAry(client)

}
